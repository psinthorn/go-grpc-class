package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/psinthorn/go-grpc-class/services/greet/greetpb"
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

	//doUnary(client)
	doServerStreaming(client)

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

// Server streaming
func doServerStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Start Unary client to Server RPC...")

	// Create greet manytimes paylaod
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sinthorn",
			LastName:  "Pradutnam",
		},
	}

	resStream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error call Greet Server Stream RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		// if we reached end of stream then break the service
		if err == io.EOF {
			break
		}

		// Iff error
		if err != nil {
			log.Fatalf("Error while reading Stream: %v", err)
		}

		log.Printf("Response from stream server: %v", msg.GetResults())
	}

}
