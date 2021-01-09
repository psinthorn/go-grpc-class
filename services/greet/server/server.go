package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/psinthorn/go-grpc-class/services/greet/greetpb"
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

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet Stream function was invoked with req: %v", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Results: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// read to end of file streaming from client
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Results: result,
			})
		}

		if err != nil {
			log.Fatal("error while reading streaming client file %v ", err)
		}

		firstname := req.GetGreet().GetFirstName()
		result += " Hello " + firstname

	}
}

func main() {
	// Log Unary server start
	fmt.Println("gRPC Server Start")

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
