package main

import (
	"context"
	"crud-api-mongodb/blog/blogpb"
	"crud-api-mongodb/blog/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

	return &blogpb.CreateBlogResponse{Blog: blogpb.CreateBlogFromBlogItem(blogItem)}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogId := req.GetBlogId()

	data, err := findBlogById(ctx, blogId)
	if err != nil {
		return nil, err
	}
	return &blogpb.ReadBlogResponse{Blog: blogpb.CreateBlogFromBlogItem(data)}, nil
}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogId := req.GetBlogId()
	blog, err := findBlogById(ctx, blogId)
	if err != nil {
		return nil, err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": blog.ID})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed while deleting: %v", err))
	}

	return &blogpb.DeleteBlogResponse{BlogId: blogId}, status.Errorf(codes.Unimplemented, "method DeleteBlog not implemented")
}

func findBlogById(ctx context.Context, blogId string) (*model.BlogItem, error) {
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid blog id: %v", blogId))
	}

	result := collection.FindOne(ctx, bson.M{"_id": oid})

	data := &model.BlogItem{}
	if err := result.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find blog with specified ID: %v", err))
	}
	return data, nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()
	blogId := blog.GetId()
	currBlog, err := findBlogById(ctx, blogId)
	if err != nil {
		return nil, err
	}

	currBlog.Title = blog.GetTitle()
	currBlog.AuthorID = blog.GetAuthorId()
	currBlog.Content = blog.GetContent()

	_, err = collection.ReplaceOne(ctx, bson.M{"_id": currBlog.ID}, currBlog)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Cannot update object: %v", err))
	}
	return &blogpb.UpdateBlogResponse{Blog: blogpb.CreateBlogFromBlogItem(currBlog)}, nil
}

func (*server) ListBlog(_ *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {

	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Listing blogs failed: %v\n", err))
	}
	defer func() { _ = cur.Close(context.Background()) }()

	for cur.Next(context.Background()) {
		data := &model.BlogItem{}
		if err := cur.Decode(data); err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Error while decoding data: %v", err))
		}
		stream.Send(&blogpb.ListBlogResponse{Blog: blogpb.CreateBlogFromBlogItem(data)})
	}
	return nil
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
