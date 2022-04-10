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

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) (err error) {
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

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			result := &greetpb.LongGreetResponse{Result: result}
			return stream.SendAndClose(result)

		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v \n", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result = fmt.Sprintf("%vHello %v! ", result, firstName)
	}
}
func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming BiDi request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v \n", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		response := &greetpb.GreetEveryoneResponse{Result: fmt.Sprintf("Hello %s !", firstName)}
		if err := stream.Send(response); err != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "the client canceled the request") // avoid expensive operation
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName + "!"
	res := &greetpb.GreetWithDeadlineResponse{Result: result}
	return res, nil
}

func main() {
	fmt.Println("Starting server...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	creds, sslError := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")
	if sslError != nil {
		log.Fatalf("Failed to serve: %v", sslError)
		return
	}

	s := grpc.NewServer(grpc.Creds(creds))
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve:%v\n", err)
	}
}
