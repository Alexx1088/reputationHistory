package handler

import (
	"context"
	"log"

	pb "github.com/Alexx1088/reputationhistory/proto"
)

type ReputationServer struct {
	pb.UnimplementedReputationServiceServer
}

func NewReputationServer() *ReputationServer {
	return &ReputationServer{}
}

func (s *ReputationServer) AddEntry(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error) {
	log.Printf("Received AddEntry request: user_id=%s, action=%s, score_change=%d",
		req.GetUserId(), req.GetAction(), req.GetScoreChange())

	// TODO: сохранять в БД, Redis или просто логировать (на старте можно просто печатать)

	return &pb.AddEntryResponse{Success: true}, nil
}
