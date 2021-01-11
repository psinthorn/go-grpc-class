package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

// Unary
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum service start with req: %v", req)
	num_1 := req.Num_1
	num_2 := req.Num_2
	Sum := num_1 + num_2
	res := &calculatorpb.SumResponse{
		Result: Sum,
	}
	return res, nil
}

// Client-Side Streaming
func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Println("Client-Side Streaming Server Start...")

	// รับค่าจาก client ที่ส่งแบบ streaming
	// วนลูปอ่านค่า
	for {
		var result = 0
		streamData, err := stream.Recv()

		for _, num := range streamData {
			result += num
			fmt.Println(result)
		}

		// หากอ่านค่าจากไฟล์ที่ส่งมาจนหมด จะได้ค่า error == EOF หมายถึงอ่านค่า streaming สำเร็จแล้วให้ นำค่าที่อ่านได้ไปทำตาม Business Logic และ Return ค่ากกลับให้ Client
		if err == EOF {
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error on reading streaming file %v ", err)
		}

	}

}

func main() {
	fmt.Println("Unary Server for Calculator Summary Start... ")

	// New gRPC server
	s := grpc.NewServer()

	// Register sum service server
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Creat Listener
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error to create listener on: %v", err.Error())
	}

	// Serve server
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error to create listener on: %v", err.Error())
	}
}
