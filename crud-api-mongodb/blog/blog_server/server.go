package main

import (
	"crud-api-mongodb/blog/blogpb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // show <file name>:<line number>

	fmt.Println("Starting blog_server...")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	var s = grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})
	reflection.Register(s)

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch // Block until a signal is received
	fmt.Println("Stopping the server")
	s.Stop()
	listen.Close()
}
