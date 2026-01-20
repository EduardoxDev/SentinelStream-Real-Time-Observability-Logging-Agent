package grpc
import (
	"context"
	"io"
	"log"
	"sync"
	"time"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "observability-system/proto/gen"
)
type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
	subscribers map[string][]chan *pb.MetricData
	mu          sync.RWMutex
}
func NewMetricsServer() *MetricsServer {
	return &MetricsServer{
		subscribers: make(map[string][]chan *pb.MetricData),
	}
}
func (s *MetricsServer) StreamMetrics(stream pb.MetricsService_StreamMetricsServer) error {
	for {
		metric, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Received metric from %s: CPU=%.2f%%, Memory=%.2f%%",
			metric.ContainerName, metric.CpuPercent, metric.MemoryPercent)
		s.broadcastMetric(metric)
		if err := stream.Send(&pb.MetricResponse{
			Success: true,
			Message: "Metric received",
		}); err != nil {
			return err
		}
	}
}
func (s *MetricsServer) GetContainerMetrics(ctx context.Context, req *pb.ContainerRequest) (*pb.ContainerMetricsResponse, error) {
	if req.ContainerId == "" {
		return nil, status.Error(codes.InvalidArgument, "container_id is required")
	}
	return &pb.ContainerMetricsResponse{
		Metrics: []*pb.MetricData{},
	}, nil
}
func (s *MetricsServer) GetHistoricalMetrics(ctx context.Context, req *pb.HistoricalRequest) (*pb.HistoricalResponse, error) {
	if req.ContainerId == "" {
		return nil, status.Error(codes.InvalidArgument, "container_id is required")
	}
	return &pb.HistoricalResponse{
		Metrics:    []*pb.MetricData{},
		TotalCount: 0,
	}, nil
}
func (s *MetricsServer) SubscribeToMetrics(req *pb.SubscriptionRequest, stream pb.MetricsService_SubscribeToMetricsServer) error {
	ch := make(chan *pb.MetricData, 100)
	s.mu.Lock()
	if req.AllContainers {
		s.subscribers["*"] = append(s.subscribers["*"], ch)
	} else {
		for _, containerID := range req.ContainerIds {
			s.subscribers[containerID] = append(s.subscribers[containerID], ch)
		}
	}
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		s.mu.Unlock()
		close(ch)
	}()
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case metric := <-ch:
			if err := stream.Send(metric); err != nil {
				return err
			}
		}
	}
}
func (s *MetricsServer) broadcastMetric(metric *pb.MetricData) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if subs, ok := s.subscribers[metric.ContainerId]; ok {
		for _, ch := range subs {
			select {
			case ch <- metric:
			default:
			}
		}
	}
	if subs, ok := s.subscribers["*"]; ok {
		for _, ch := range subs {
			select {
			case ch <- metric:
			default:
			}
		}
	}
}
func (s *MetricsServer) BroadcastMetricFromCollector(metric *pb.MetricData) {
	s.broadcastMetric(metric)
}