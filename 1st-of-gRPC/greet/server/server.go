package main

import (
	"fmt"
	"log"
	"net"

	"github.com/psinthorn/go-grpc-class/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Hello gRPC")
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen on : %v", err)
	}

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to listen on : %v", err)
	}
}
