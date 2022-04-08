package main

import (
	"client-server-boilerplate/greet/greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func main() {
	fmt.Println("Starting a client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v\n\n", err)
	}

	defer func() { _ = cc.Close() }()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created client %f", c)

	//doUnary(c)

	doServerStreaming(err, c)

}

func doServerStreaming(err error, c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{FirstName: "Hiago", LastName: "Oliveira"},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		fmt.Printf("Received message: %v\n", msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{FirstName: "Hiago", LastName: "Oliveira"},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v\n", res.Result)
}
