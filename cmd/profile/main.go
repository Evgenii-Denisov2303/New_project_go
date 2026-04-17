package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"inkflow/internal/config"
	"inkflow/internal/storage/postgres"
	profilegrpc "inkflow/internal/transport/grpc/profile"
	profilev1 "inkflow/proto/profile/v1"
)

func main() {
	cfg := config.MustLoad()

	storage, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("create postgres storage: %v", err)
	}
	defer storage.Close()

	if err := storage.Ping(); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+cfg.ProfileGRPCPort)
	if err != nil {
		log.Fatalf("listen tcp: %v", err)
	}

	grpcServer := grpc.NewServer()
	profileServer := profilegrpc.New(storage)

	profilev1.RegisterProfileServiceServer(grpcServer, profileServer)

	log.Printf("profile grpc service starting... env=%s port=%s", cfg.Env, cfg.ProfileGRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve grpc: %v", err)
	}
}
