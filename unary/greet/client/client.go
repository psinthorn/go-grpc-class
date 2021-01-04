package main

import (
	"context"
	"fmt"
	"log"

	"github.com/psinthorn/go-grpc-class/unary/greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Client start to Dial to Unary server")
	// Create client connect to server payload
	ccp, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can't Connect To : %v", err)
	}

	// Call Service on server with connection payload
	client := greetpb.NewGreetServiceClient(ccp)

	doUnary(client)

}

func doUnary(client greetpb.GreetServiceClient) {

	// Create req with payload
	fmt.Println("Start Unary client to Server RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sinhtorn",
			LastName:  "Pr.",
		},
	}

	// Get payload response from server and display
	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error call Greet RPC: %v", err)
	}

	log.Fatalf("Response from server: %v", res.Results)
}
