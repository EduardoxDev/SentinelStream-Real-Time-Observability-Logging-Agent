package main
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"observability-system/internal/infrastructure/adapters"
	"observability-system/internal/storage"
	ws "observability-system/internal/websocket"
)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
type Server struct {
	hub         *ws.Hub
	influxDB    *storage.InfluxDBStorage
	redisClient *redis.Client
}
func main() {
	log.Println("🚀 Starting Observability Server...")

	influxDB := storage.NewInfluxDBStorage(
		getEnv("INFLUXDB_URL", "http://localhost:8086"),
		getEnv("INFLUXDB_TOKEN", "my-super-secret-token"),
		getEnv("INFLUXDB_ORG", "observability"),
		getEnv("INFLUXDB_BUCKET", "metrics"),
	)
	defer influxDB.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})
	defer redisClient.Close()
	hub := ws.NewHub()
	go hub.Run()
	server := &Server{
		hub:         hub,
		influxDB:    influxDB,
		redisClient: redisClient,
	}
	go server.broadcastMetrics()

	http.HandleFunc("/ws", server.handleWebSocket)
	http.HandleFunc("/api/containers", server.handleContainers)
	http.HandleFunc("/api/metrics", server.handleMetrics)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	port := getEnv("PORT", "8080")
	log.Printf("🌐 Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &ws.Client{
		Hub:  s.hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	s.hub.Register(client)
	go client.WritePump()
	go client.ReadPump()
}
func (s *Server) handleContainers(w http.ResponseWriter, r *http.Request) {
	processes := []map[string]string{
		{"ID": "process-observability-agent", "Names": "observability-agent"},
		{"ID": "process-observability-server", "Names": "observability-server"},
		{"ID": "process-redis-server", "Names": "redis-server"},
		{"ID": "process-influxdb", "Names": "influxdb"},
		{"ID": "process-chrome", "Names": "chrome"},
		{"ID": "process-vscode", "Names": "vscode"},
		{"ID": "process-powershell", "Names": "powershell"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processes)
}
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	containerID := r.URL.Query().Get("container_id")
	if containerID == "" {
		http.Error(w, "container_id required", http.StatusBadRequest)
		return
	}
	metrics, err := s.influxDB.QueryMetrics(r.Context(), containerID, 1*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
func (s *Server) broadcastMetrics() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	processCollector, _ := adapters.NewProcessCollectorAdapter()
	defer processCollector.Close()
	for range ticker.C {
		ctx := context.Background()
		processes, err := processCollector.ListContainers(ctx)
		if err != nil {
			continue
		}
		type MetricData struct {
			ContainerID   string    `json:"ContainerID"`
			ContainerName string    `json:"ContainerName"`
			CPUPercent    float64   `json:"CPUPercent"`
			MemoryUsage   uint64    `json:"MemoryUsage"`
			MemoryLimit   uint64    `json:"MemoryLimit"`
			MemoryPercent float64   `json:"MemoryPercent"`
			NetworkRx     uint64    `json:"NetworkRx"`
			NetworkTx     uint64    `json:"NetworkTx"`
			Timestamp     time.Time `json:"Timestamp"`
		}
		var allStats []MetricData
		for _, processName := range processes {
			metrics, err := processCollector.CollectMetrics(ctx, processName)
			if err != nil || metrics == nil {
				continue
			}
			stats := MetricData{
				ContainerID:   metrics.ContainerID,
				ContainerName: metrics.ContainerName,
				CPUPercent:    metrics.CPUPercent,
				MemoryUsage:   metrics.MemoryUsage,
				MemoryLimit:   metrics.MemoryLimit,
				MemoryPercent: metrics.MemoryPercent,
				NetworkRx:     metrics.NetworkRx,
				NetworkTx:     metrics.NetworkTx,
				Timestamp:     metrics.Timestamp,
			}
			allStats = append(allStats, stats)
		}
		s.hub.BroadcastMetrics(allStats)
	}
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}