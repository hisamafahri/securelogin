package main

import (
	"log"

	"github.com/hisamafahri/securelogin/infrastructure/pgsql"
)

func main() {
	_, err := pgsql.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if err := pgsql.Close(); err != nil {
			log.Printf("failed to close database connection: %v", err)
		}
	}()

	r := NewServer()

	addr := ":8080"
	log.Printf("Server is running on http://localhost%s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("[error]: server %v", err)
	}
}
