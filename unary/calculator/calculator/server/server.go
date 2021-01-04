package main

import (
	"context"
	"fmt"
	"log"
	"net"

	calculatorpb "github.com/psinthorn/go-grpc-class/calculator/proto"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum service start with req: %v", req)
	num_1 := req.Num_1
	num_2 := req.Num_2
	Sum := num_1 + num_2
	res := &calculatorpb.SumResponse{
		Result: Sum,
	}
	return res, nil
}

func main() {
	fmt.Println("Unary Server for Calculator Summary Start... ")

	// New gRPC server
	s := grpc.NewServer()

	// Register sum service server
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Creat Listener
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error to create listener on: %v", err.Error())
	}

	// Serve server
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error to create listener on: %v", err.Error())
	}
}
