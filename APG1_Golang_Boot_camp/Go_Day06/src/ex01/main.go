package main

import (
	"blog/server"
	"blog/db"
	"log"
	"context"
)

func main() {
	conn, err := db.NewPostgres("krillkovalev", "108814", "localhost", "5432", "my_blog")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer conn.Close(context.Background())
}
