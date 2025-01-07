package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"military/military"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Anomaly struct {
	Session_ID string
	Timestamp string
	Frequency float64
}

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

	db := Init()

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
		session_id := msg.GetSessionId()
		timestamp := msg.GetTimestamp()
		sum += frequency 
		mean = sum / float64(count)
		sq = (frequency - mean) * (frequency - mean) 
		sumSq += sq
		dispersion = sumSq / float64(count)
		std = math.Sqrt(dispersion)
		data := Anomaly {
			Session_ID: 					session_id,
			Timestamp: 						timestamp,
			Frequency:						frequency,							
		}
		log.Printf("Session_ID: %s, Timestamp: %s Frequency: %f, Mean: %f, Standart Deviation: %f, Count of values procceded: %d", session_id, timestamp, frequency, mean, std, count)
		// log.Printf("Received message: %+v", msg)

		if math.Abs(frequency - mean) > *flag_k*std {
			db.Table("anomalies")
			if result := db.Create(&data); result.Error != nil {
				fmt.Println(result.Error)
			}
			log.Printf("Anomaly Detected! Session_ID: %s, Timestamp: %s, Frequency: %f, Mean: %f, Standart Deviation: %f", session_id, timestamp, frequency, mean, std)
		}
	}

}