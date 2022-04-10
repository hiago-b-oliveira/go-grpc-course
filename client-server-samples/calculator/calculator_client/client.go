package main

import (
	"client-server-boilerplate/calculator/calculatorpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
)

func main() {
	fmt.Println("Starting a client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v\n\n", err)
	}

	defer func() { _ = cc.Close() }()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//fmt.Printf("Created client %f", c)

	//doUnarySum(c)
	//doDiBiFindMax(c)

	calcSquareGRPC(c, 16)
	calcSquareGRPC(c, -16)
}

func calcSquareGRPC(c calculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
		}
		return
	}
	fmt.Printf("Result of square root of %v: %v\n", number, res.GetNumberRoot())
}

func doDiBiFindMax(c calculatorpb.CalculatorServiceClient) {
	stream, _ := c.FindMax(context.Background())
	waitc := make(chan struct{})

	maxResults := make([]int32, 0)
	go func() {
		inputs := []int32{1, 5, 3, 6, 2, 20}
		for _, input := range inputs {
			stream.Send(&calculatorpb.FindMaxRequest{Value: input})
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			maxResults = append(maxResults, res.GetValue())
		}
	}()
	<-waitc
	fmt.Printf("Max Results: %v", maxResults)
}

func doUnarySum(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.CalculatorRequest{
		Calculation: &calculatorpb.Calculation{A: 1, B: 2},
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Sum: %v\n", res.Result)
}
