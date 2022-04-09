package main

import (
	"client-server-boilerplate/calculator/calculatorpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(_ context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	a, b := req.Calculation.A, req.Calculation.B
	sum := a + b
	resp := &calculatorpb.CalculatorResponse{
		Result: sum,
	}
	return resp, nil
}

func (*server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {
	max := int32(-1)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if req.GetValue() > max {
			max = req.GetValue()
			stream.Send(&calculatorpb.FindMaxResponse{Value: max})
		}
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
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve:%v\n", err)
	}
}
