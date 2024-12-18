package main

import (
	"context"
	"io"
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

	for {
		msg, err := r.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error receiving data: %v", err)
		}

		log.Printf("Received message: %+v", msg)
	}

}