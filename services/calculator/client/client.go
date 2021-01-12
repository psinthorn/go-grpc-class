package main

import (
	"context"
	"fmt"
	"log"

	calculatorpb "github.com/psinthorn/go-grpc-class/services/calculator/proto"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator client start to connect to server...")

	// Create client connection to server paylaod
	ccp, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed connect to server: %v", err)
	}

	// Call service on server with client connection payload
	client := calculatorpb.NewCalculatorServiceClient(ccp)

	// Unary client
	// doUnarySum(client)

	// Client of Client-Side streaming
	doClientSideAverage(client)

}

func doUnarySum(client calculatorpb.CalculatorServiceClient) {
	// Create request paylaod
	req := &calculatorpb.SumRequest{
		Num_1: 3,
		Num_2: 2,
	}

	// Response with result payload from server
	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error on sum: %v", res)
	}

	log.Fatalf("Sum of req is %v", res)
}

func doClientSideAverage(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Client-Side Streaming start...")

	numbers := []int32{1, 6, 5, 7, 9, 11}

	stream, err := client.Average(context.Background())
	for _, num := range numbers {
		fmt.Printf("sending number: %v \n", num)
		stream.Send(&calculatorpb.AverageRequest{
			Num: num,
		})
		if err != nil {
			log.Fatalf("Error on streaming data to server %v\n", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("Error on receiving %v ", err)
	}

	fmt.Printf("Average is: %v", res.GetResult())

}
