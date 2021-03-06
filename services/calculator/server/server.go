package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	calculatorpb "github.com/psinthorn/go-grpc-class/services/calculator/proto"
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
	var (
		num_count int32 = 0
		sum       int32 = 0
	)
	// รับค่าจาก client ที่ส่งแบบ streaming
	// วนลูปอ่านค่า
	for {
		streamData, err := stream.Recv()
		// หากอ่านค่าจากไฟล์ที่ส่งมาจนหมด จะได้ค่า error == EOF หมายถึงอ่านค่า streaming สำเร็จแล้วให้ นำค่าที่อ่านได้ไปทำตาม Business Logic และ Return ค่ากกลับให้ Client
		if err == io.EOF {
			average := sum / num_count
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: average,
			})
		}

		if err != nil {
			log.Fatalf("Error on reading streaming file %v ", err)
		}

		sum += streamData.GetNum()
		num_count++
		fmt.Printf("Num Count is: %v \n", num_count)
		fmt.Printf("Sum round %v is %v ", num_count, sum)

	}

}

func (*server) Maximum(stream calculatorpb.CalculatorService_MaximumServer) error {
	// Receiving data from client request
	// รับข้อมูลจาก client request
	maxNum := int32(0)
	for {
		streamData, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("Reading data completed... %v ", err)
			return nil
		}
		if err != nil {
			fmt.Printf("Error while reading req data stream %v ", err)
			return err
		}

		reqNum := streamData.GetNum()
		fmt.Printf("Current number is: %v \n", maxNum)
		if reqNum > maxNum {
			maxNum = reqNum
			fmt.Printf("Maximum is: %v \n", maxNum)

			result := &calculatorpb.MaximumResponse{
				Result: maxNum,
			}
			stream.Send(result)
		}
	}

}

func main() {
	fmt.Println("Server for Calculator is start... ")

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
