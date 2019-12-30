package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ArchieSpinos/grpc-go/calculator/calcpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) SumNums(ctx context.Context, in *calcpb.CalcRequest) (*calcpb.CalcResponse, error) {
	return &calcpb.CalcResponse{Result: in.Numbers.GetNum1() + in.Numbers.GetNum2()}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen to %v", err)
	}
	grpcServer := grpc.NewServer()
	calcpb.RegisterCalcServiceServer(grpcServer, &server{})

	//Start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server: %s", err)
	}
}
