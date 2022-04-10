package main

import (
	"context"
	"crud-api-mongodb/blog/blogpb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	fmt.Println("Staring Blog Client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v\n", err)
	}
	defer func() { _ = cc.Close() }()

	c := blogpb.NewBlogServiceClient(cc)

	blog := &blogpb.Blog{
		Id:       "",
		AuthorId: "Hiago",
		Title:    "Go Course",
		Content:  "This a blog about a Go Course",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		fmt.Printf("Creating a blog failed: %v", err)
		return
	}
	fmt.Printf("Blog created: %v\n", res)
}
