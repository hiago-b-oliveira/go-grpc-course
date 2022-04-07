package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"validating-configs/greet/greetpb"
)

func main() {
	fmt.Println("Starting a client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v\n\n", err)
	}

	defer func() { _ = cc.Close() }()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Created client %f", c)
}
