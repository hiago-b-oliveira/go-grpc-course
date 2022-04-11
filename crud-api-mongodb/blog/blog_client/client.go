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

	//blogCreated := createBlog(c, &blogpb.Blog{
	//	Id:       "",
	//	AuthorId: "Hiago",
	//	Title:    "Go Course",
	//	Content:  "This a blog about a Go Course",
	//})
	//
	//blog1, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogCreated.GetId()})
	//fmt.Printf("Reading blog... Blog: %v, Error: %v\n", blog1, err)
	//
	//blog2, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "123-invalid-id"})
	//fmt.Printf("Reading blog... Blog: %v, Error: %v\n", blog2, err)
	//
	//blog3, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "ffffffffffffffffffffffff"})
	//fmt.Printf("Reading blog... Blog: %v, Error: %v\n", blog3, err)
	//
	//blogCreated.Title = "gRPC Go Course"
	//updatedBlog, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: blogCreated})
	//fmt.Printf("Updated blog: %v, err: %v\n", updatedBlog, err)

	_, err = c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: "6252f05598c9bd7a428a9a16"})
	if err != nil {
		fmt.Printf("Delete failed: %v", err)
	}

}

func createBlog(c blogpb.BlogServiceClient, blog *blogpb.Blog) *blogpb.Blog {
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		fmt.Printf("Creating a blog failed: %v", err)
		return nil
	}
	fmt.Printf("Blog created: %v\n", res)
	return res.GetBlog()
}
