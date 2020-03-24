package main

import (
	"fmt"
	"log"
	"os"

	users "github.com/microapis/users-api/run"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatalln("missing env variable HOST")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("missing env variable PORT")
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		log.Fatalln("missing env variable POSTGRES_DSN")
	}

	users.Run(addr, postgresDSN)
}
