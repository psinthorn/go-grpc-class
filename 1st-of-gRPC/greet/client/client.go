package main

import (
	"fmt"
	"log"

	"github.com/psinthorn/go-grpc-class/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Hello I'm gRPC Client")

	clientConect, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server : %v", err)
	}

	defer clientConect.Close()

	client := greetpb.NewGreetServiceClient(clientConect)
	fmt.Printf("Client created %f", client)

}
