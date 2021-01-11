package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// doUnary(client)
	//doServerStreaming(client)
	doClientStreaming(client)

}

func doUnary(client greetpb.GreetServiceClient) {

	// Create req with payload
	fmt.Println("Start Unary RPC...")
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
	fmt.Println("Start Server streaming client RPC...")

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

// client-side streaming
func doClientStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Start streaming client RPC...")

	// Prepare data or slice of data
	dataRequest := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sinthorn",
				LastName:  "Pr.",
			},
		},

		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nut",
				LastName:  "Pr.",
			},
		},

		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Na Phansa",
				LastName:  "Pr.",
			},
		},

		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ichi",
				LastName:  "Pr.",
			},
		},

		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Taro",
				LastName:  "Pr.",
			},
		},
	}
	// create stream by calling
	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet %v\n", err)
	}
	// iterate slice and send data to server
	for _, data := range dataRequest {
		fmt.Printf("start to streaming data to serve %v \n", data)
		stream.Send(data)
		time.Sleep(100 * time.Millisecond)
	}
	// close and receive results from server
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error on close and receive %v \n", err)
	}

	fmt.Printf("Results: %v \n", res)
}
