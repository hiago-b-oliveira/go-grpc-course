package main

import (
	"client-server-boilerplate/greet/greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(_ context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := greetpb.GreetResponse{
		Result: result,
	}

	return &res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) (err error) {
	fmt.Printf("GrGreetManyTimes function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		content := fmt.Sprintf("Hello %v number %v", firstName, i)
		res := &greetpb.GreetManyTimesResponse{
			Result: content,
		}
		err = stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func main() {
	fmt.Println("Starting server...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve:%v\n", err)
	}
}
