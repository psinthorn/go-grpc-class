package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/psinthorn/go-grpc-class/services/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

// Unary
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

// Server-Side Streaming
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet Server Streaming Server function was invoked with req: %v", req)
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

// Client-Side Streaming
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("Greet Client Streaming Server function was invoked\n")
	result := ""

	// เนื่องเป็น Client Streaming ที่มีการส่งข้อมูลแบบต่อเนื่องมาให้ที่ Server ดังนั้น ขั้นตอนการทำงานของฝั่ง Server จึงเป็นดังนี้
	// 1. วนลูปอ่านข้อมูลที่ที่ส่งมาจากฝั่ง client ลงในตัวแปรด้วยเมธอด stream.Recv() โดยที่เมธอด Recv() คืนค่ามาสองตัวแปรคือค่าที่ส่งมาจากฝั่ง client และค่า error
	// 2. ตรวจสอบ error ที่ได้รับกลับมาว่า
	// 2.1 หาก error == io.EOF หมายความว่าอ่านค่าจากไฟล์ที่ส่งเข้ามาสำเร็จและสิ้นสุดแล้ว นำข้อมูลที่ได้ไปดำเนินการตาม Business Logic ของเรา
	// 2.2 หาก error != nil หมายความว่ามีปัญหากับข้อมูลที่ทางฝั้ง client ส่งมา ให้ return ค่า error กลับไปให้ทางฝั่ง client ได้รับทราบ
	//และ return เพื่อคืนค่าที่ดำเนินการเรียบร้อยแล้วกลับไปให้ฝั่ง Client ตามรูปแบบ model ของข้อมูลที่กำหนดใน .proto ไฟล์ ด้วยเมธอด stream.SendAndClose

	for {
		// 1. วนลูปอ่านข้อมูลที่ที่ส่งมาจากฝั่ง client ลงในตัวแปรด้วยเมธอด stream.Recv()
		req, err := stream.Recv()

		if err == io.EOF {
			// 2.1 อ่านค่าจากไฟล์ที่ส่งเข้ามาสำเร็จและสิ้นสุดแล้ว นำข้อมูลที่ได้ไปดำเนินการตาม Business Logic ของเรา และ return กลับไปให้ทางฝั่ง client
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Results: result,
			})
		}

		// 2.1 หากมีข้อผิดพลาด ให้ return ค่า error กลับไปให้ทางฝั่ง client ได้รับทราบ
		if err != nil {
			log.Fatal("error while reading streaming client file %v ", err)
		}

		// 2.2 นำข้อมูลที่ได้รับมาดำเนินการตาม Business Logic (หากข้อมูลที่ส่งมายังไม่หมดจะดำเนินการตาม Logic จนกว่าจะ EOF)
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		result += " Hello " + firstName + " " + lastName

	}
}

func main() {

	// Log Unary server start
	fmt.Println("gRPC Server Start")

	// Server port env
	envPort := os.Getenv("PORT")

	if envPort == "" {
		envPort = "50051"
	}

	fmt.Printf("Env server port is: %s", envPort)

	// Need grpc server
	s := grpc.NewServer()

	// Register server serivce
	greetpb.RegisterGreetServiceServer(s, &server{})

	// Need listener
	listener, err := net.Listen("tcp", "0.0.0.0:"+envPort)
	if err != nil {
		log.Fatalf("Failed to connect to : %v", err)
	}

	// Serve server
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to connect to : %v", err)
	}

}
