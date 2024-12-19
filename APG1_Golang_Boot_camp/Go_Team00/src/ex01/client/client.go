package main

import (
	"context"
	"io"
	"log"
	"military/military"
	"flag"
	"math"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const k float64 = 2.0

func main() {
	flag_k := flag.Float64("k", k, "STD anomaly coefficient")
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	c := military.NewDeviceServiceClient(conn)
	
	ctx := context.Background()

	r, err := c.StreamData(ctx, &military.ConnectionRequest{})
	if err != nil {
		log.Fatalf("error calling request Stream data: %v", err)
	}

	sum := 0.
	sumSq := 0.
	count := 0
	mean := 0.
	frequency := 0.
	sq := 0.
	dispersion := 0.
	std := 0.

	for {
		msg, err := r.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error receiving data: %v", err)
		}
		count++
		frequency = msg.GetFrequency()
		sum += frequency 
		mean = sum / float64(count)
		sq = (frequency - mean) * (frequency - mean) 
		sumSq += sq
		dispersion = sumSq / float64(count)
		std = math.Sqrt(dispersion)

		log.Printf("Frequency: %f, Mean: %f, Standart Deviation: %f, Count of values procceded: %d", frequency, mean, std, count)
		// log.Printf("Received message: %+v", msg)

		if math.Abs(frequency - mean) > *flag_k*std {
			log.Printf("Anomaly Detected! Frequency: %f, Mean: %f, Standart Deviation: %f", frequency, mean, std)
		}
	}

}