package main

import (
	"database/sql"
	"log"
	"os"

	ticketServer "ticket/cmd"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := ticketServer.NewTicketServer(db)
	go func() {
		log.Fatal(ticketServer.StartGRPCServer(server, ":9090"))
	}()
	log.Fatal(ticketServer.StartHTTPServer(server, ":8080"))
}
