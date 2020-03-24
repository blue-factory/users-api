package main

import (
	"log"
	"os"

	users "github.com/microapis/users-api/run"
)

func main() {
	port := os.Getenv("PORT")
	postgresDSN := os.Getenv("POSTGRES_DSN")

	if port == "" {
		log.Fatalln("missing env variable PORT")
	}

	if postgresDSN == "" {
		log.Fatalln("missing env variable POSTGRES_DSN")
	}

	users.Run(port, postgresDSN)
}
