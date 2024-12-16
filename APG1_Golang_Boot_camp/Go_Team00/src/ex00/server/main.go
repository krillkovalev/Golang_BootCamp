package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"math/rand"
	"military/military"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"time"
)

type server struct {
	military.UnimplementedDeviceServiceServer
	res map[string]*military.DeviceData
}

func (s *server) StreamData(ctx context.Context, req *military.ConnectionRequest) {
	mean := rand.Intn(-10+1-10) + 10
	std := 0.3 + rand.Float64() * (1.5 - 0.3)
	frequency := float64(mean) + std * rand.NormFloat64()
	timestamp := time.Now().UTC().String()
	data := military.DeviceData{
		SessionId: 			uuid.NewString(),
		Frequency: 			frequency,
		Timestamp: 			timestamp,
	} 
	


}