package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"validating-configs/greet/greetpb"
)

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &greetpb.UnimplementedGreetServiceServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve:%v\n", err)
	}
}
