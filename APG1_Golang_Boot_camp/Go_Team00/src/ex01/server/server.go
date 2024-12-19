package main

import (
	"log"
	"math/rand"
	"military/military"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	military.UnimplementedDeviceServiceServer
}

func (s *server) StreamData(req *military.ConnectionRequest, stream grpc.ServerStreamingServer[military.DeviceData]) error{
	timestamp := time.Now().UTC().String()
	
	for {
		rand.New(rand.NewSource(time.Now().UnixNano()))
		min := -10.0
		max := 10.0
		mean := min + rand.Float64() * (max - min)
		std := 0.3 + rand.Float64() * (1.5 - 0.3)
		frequency := float64(mean) + std * rand.NormFloat64()
		data := military.DeviceData{
			SessionId: 			uuid.NewString(),
			Frequency: 			frequency,
			Timestamp: 			timestamp,
		}
		if err := stream.Send(&data); err != nil {
			return err
		}

		time.Sleep(100 * time.Millisecond)
	}
	

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen 50051 port: %v", err)
	}

	s := grpc.NewServer()
	military.RegisterDeviceServiceServer(s, &server{})
	log.Printf("grpc server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}