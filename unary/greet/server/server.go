package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/psinthorn/go-grpc-class/unary/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with req: %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName

	res := &greetpb.GreetResponse{
		Results: result,
	}
	return res, nil
}

func main() {
	// Log Unary server start
	fmt.Println("Unary server start")

	// Need grpc server
	s := grpc.NewServer()

	// Register server serivce
	greetpb.RegisterGreetServiceServer(s, &server{})

	// Need listener
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to connect to : %v", err)
	}

	// Serve server
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to connect to : %v", err)
	}
}
