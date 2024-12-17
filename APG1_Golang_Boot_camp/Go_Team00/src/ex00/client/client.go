package main

import (
	"context"
	"log"
	"military/military"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	c := military.NewDeviceServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.StreamData(ctx, &military.ConnectionRequest{})
	if err != nil {
		log.Fatalf("error calling request Stream data: %v", err)
	}

	log.Printf("Response from server: %s", r)
}