package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Alexx1088/reputationhistory/internal/consumer"
	"github.com/Alexx1088/reputationhistory/internal/handler"
	pb "github.com/Alexx1088/reputationhistory/proto"
	_ "github.com/lib/pq"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("ReputationHistory service starting...")

	// подключение к БД
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	// поднимаем gRPC сервер
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterReputationServiceServer(grpcServer, handler.NewReputationServer())

	// запускаем consumer в отдельной горутине
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		c := consumer.NewReputationConsumer(db)
		if err := c.Run(ctx); err != nil {
			log.Fatalf("consumer error: %v", err)
		}
	}()

	fmt.Println("ReputationHistory service is listening on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
