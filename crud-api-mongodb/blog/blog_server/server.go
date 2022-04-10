package main

import (
	"context"
	"crud-api-mongodb/blog/blogpb"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
)

var collection *mongo.Collection

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	blogItem := blog.AsBlogItem()

	res, err := collection.InsertOne(ctx, blogItem)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)
	blogItem.ID = oid

	createBlogResponse := blogpb.CreateBlogResponseFromBlogItem(blogItem)
	return createBlogResponse, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // show <file name>:<line number>

	grpcServer, tcpListener := createGrpcServer()
	mongoClient := createMongoClient()

	collection = mongoClient.Database("mydb").Collection("blog")

	defer func() {
		fmt.Println("Disconnecting from mongo")
		_ = mongoClient.Disconnect(context.Background())

		fmt.Println("Stopping the server")
		grpcServer.Stop()
		tcpListener.Close()
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch // Block until a signal is received
}

func createGrpcServer() (*grpc.Server, net.Listener) {
	fmt.Println("Starting blog gRPC server...")
	tcpListener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to tcpListener: %v\n", err)
	}

	grpcServer := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(tcpListener); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()
	return grpcServer, tcpListener
}

func createMongoClient() *mongo.Client {
	clientOpts := options.Client().ApplyURI("mongodb://admin:password@localhost:27017")
	mongoClient, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatalf("Mongo connection failed: %v\n", err)
	}
	if err := mongoClient.Connect(context.Background()); err != nil {
		log.Fatalf("Mongo connection failed: %v\n", err)
	}
	fmt.Println("Mongodb connected")
	return mongoClient
}
