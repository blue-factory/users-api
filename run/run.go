package run

import (
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/microapis/users-api/database"
	pb "github.com/microapis/users-api/proto"
	usersSvc "github.com/microapis/users-api/rpc/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run ...
func Run(port string, postgresDSN string) {
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
