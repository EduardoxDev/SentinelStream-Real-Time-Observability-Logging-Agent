package grpc
import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "observability-system/proto/gen"
)
type MetricsClient struct {
	conn   *grpc.ClientConn
	client pb.MetricsServiceClient
}
func NewMetricsClient(serverAddr string) (*MetricsClient, error) {
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}
	return &MetricsClient{
		conn:   conn,
		client: pb.NewMetricsServiceClient(conn),
	}, nil
}
func (c *MetricsClient) SendMetric(ctx context.Context, metric *pb.MetricData) error {
	stream, err := c.client.StreamMetrics(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(metric); err != nil {
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	if !resp.Success {
		log.Printf("Failed to send metric: %s", resp.Message)
	}
	return stream.CloseSend()
}
func (c *MetricsClient) SubscribeToAllMetrics(ctx context.Context) (<-chan *pb.MetricData, error) {
	stream, err := c.client.SubscribeToMetrics(ctx, &pb.SubscriptionRequest{
		AllContainers: true,
	})
	if err != nil {
		return nil, err
	}
	ch := make(chan *pb.MetricData, 100)
	go func() {
		defer close(ch)
		for {
			metric, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving metric: %v", err)
				return
			}
			ch <- metric
		}
	}()
	return ch, nil
}
func (c *MetricsClient) Close() error {
	return c.conn.Close()
}