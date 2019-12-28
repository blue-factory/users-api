package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/microapis/users-api/database"
	usersSvc "github.com/microapis/users-api/rpc/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/microapis/users-api/proto"
	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	postgresDSN := os.Getenv("POSTGRES_DSN")

	if port == "" {
		port = "5020"
		log.Println("missing env variable PORT, using default value...")
	}

	if postgresDSN == "" {
		postgresDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
		fmt.Println("missing env variable POSTGRES_DSN, using default value")
	}

	pgSvc, err := database.NewPostgres(postgresDSN)
	if err != nil {
		log.Println("PG DSN ", postgresDSN)
		log.Fatalf("Failed connect to postgres: %v", err)
	}

	srv := grpc.NewServer()
	svc := usersSvc.New(pgSvc)

	pb.RegisterUserServiceServer(srv, svc)
	reflection.Register(srv)

	log.Println("Starting Users service...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("Failed to list: %v", err)
	}

	log.Println(fmt.Sprintf("Users service, Listening on: %v", port))

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Fatal to serve: %v", err)
	}
}
