package main

import (
	"context"
	"fmt"
	"log"

	calculatorpb "github.com/psinthorn/go-grpc-class/unary/calculator/proto"

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

	doUnarySum(client)

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
