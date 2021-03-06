package main

import (
	"client-server-boilerplate/calculator/calculatorpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
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

func (*server) SquareRoot(_ context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative number: %v", number))
	}
	return &calculatorpb.SquareRootResponse{NumberRoot: math.Sqrt(float64(number))}, nil
}

func main() {
	fmt.Println("Starting server...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s) // Allow us to use a CLI: https://github.com/ktr0731/evans

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve:%v\n", err)
	}
}
