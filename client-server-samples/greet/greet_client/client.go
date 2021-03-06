package main

import (
	"client-server-boilerplate/greet/greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Starting a client...")
	creds, sslError := credentials.NewClientTLSFromFile("ssl/ca.crt", "")
	if sslError != nil {
		log.Fatalf("Error while loading CA trust certificate: %v\n", sslError)
	}
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("could not connect: %v\n\n", err)
	}

	defer func() { _ = cc.Close() }()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Println("Client Created!")

	//doUnary(c)
	//doServerStreaming(err, c)
	//doClientStreaming(err, c)

	//doBiDiStreaming(err, c)

	doUnaryGreetWithDeadline(c, 5*time.Second) // should complete
	doUnaryGreetWithDeadline(c, 1*time.Second) // should timeout
}

func doUnaryGreetWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	greeting := greetpb.Greeting{FirstName: "Hiago"}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, &greetpb.GreetWithDeadlineRequest{Greeting: &greeting})
	if err != nil {
		if statusErr, ok := status.FromError(err); ok && statusErr.Code() == codes.DeadlineExceeded {
			fmt.Println("Timeout was hit! Deadline was exceeded")
			return
		}
		log.Fatalf("Error while calling GreetWithDeadline RPC: %v\n", err)
		return
	}
	log.Printf("Response from Greet: %v", res.GetResult())
}

func doBiDiStreaming(err error, c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		names := []string{"Hiago", "John", "Patric", "Marie"}
		for _, name := range names {
			fmt.Println("Sending a request: " + name)
			req := &greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{FirstName: name},
			}
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				close(waitc)
				log.Fatalf("Error while receiving: %v", err)
				return
			}
			fmt.Printf("Received: %v\n", res)
		}
	}()
	<-waitc
}

func doClientStreaming(err error, c greetpb.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v/n", err)
	}
	for i := 0; i < 10; i++ {
		stream.Send(&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: fmt.Sprintf("%v [%v]", "Hiago", i)},
		})
	}

	result, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving the result: %v/n", err)
	}
	fmt.Printf("Result Received: %v", result.GetResult())
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
