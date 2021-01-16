package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doClientSideAverage(client)
	doFindMaximum(client)

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

	numbers := []int32{1, 6, 10, 7, 9, 110}

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

// Bidirectional streaming client
func doFindMaximum(client calculatorpb.CalculatorServiceClient) {
	// Preparing data
	numbers := []int32{1, 3, 2, 4, 5, 9, 2, 8, 10, 5, 11}

	// create client Maximum stream
	stream, err := client.Maximum(context.Background())
	if err != nil {
		fmt.Printf("Can't create Maximum stream %v ", err)
		return
	}

	wg := make(chan struct{})
	// Send data
	go func() {
		for _, req := range numbers {
			fmt.Printf("Sending number %v \n", req)
			if err := stream.Send(&calculatorpb.MaximimumRequest{Num: req}); err != nil {
				fmt.Printf("Error while sending stream to server %v ", err)
				return
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receive data
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Receiving data completed")
				return
			}

			if err != nil {
				fmt.Printf("Error while reading data %v \n ", err)
				return
			}

			result := res.GetResult()
			fmt.Printf("Maximum numbers is: %v \n", result)
		}
		close(wg)
	}()
	<-wg

}
