package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-go-course/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

var collection *mongo.Collection

type blogItem struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthorID string             `json:"author_id" bson:"author_id"`
	Content  string             `json:"content" bson:"content"`
	Title    string             `json:"title" bson:"title"`
}

func main() {
	//If we crash the go code we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	s := grpc.NewServer()
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	blogpb.RegisterBlogServiceServer(s, &server{})

	fmt.Println("Connecting to mongodb")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@localhost:27017/admin"))
	if err != nil {
		log.Fatalf("Error setting up new mongo client: [ %v ]", err)
		return
	}
	fmt.Println("Connected to mongodb")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connecting into client: [ %v ]", err)
		return
	}

	collection = client.Database("course-db").Collection("blog")

	go func() {
		fmt.Println("[ blog ] Starting server up!")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to server: %v", err)
		}
	}()

	// Waits for control + c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//    Block until signal is recieved
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Stopping the listener")
	err = lis.Close()
	if err != nil {
		log.Fatalf("Error closing listener: [ %v ]", err)
	}

	fmt.Println("Closing mongodb")
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			fmt.Printf("Error closing mongodb connection: [ %v ]", err)
		}
	}(client, ctx)

	fmt.Println("Graceful shutdown done")

}

func (*server) CreateBlog(_ context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.AuthorId,
		Content:  blog.Content,
		Title:    blog.Title,
	}

	resInput, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: [ %v ]", err),
		)
	}

	oid, ok := resInput.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Couldn't convert id: [ %v ]", err),
		)
	}

	res := &blogpb.CreateBlogResponse{Blog: &blogpb.Blog{
		Id:       oid.Hex(),
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}}

	return res, nil
}

func (*server) ReadBlog(_ context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "[ blog ] couldn't parse invalid id: [ %v ]", oid)
	}

	// create an empty struct
	data := &blogItem{}
	filter := bson.M{"_id": oid}

	result := collection.FindOne(context.Background(), filter)
	if err := result.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, "[ blog ] couldn't find given ID: [ %v ]", err)
	}

	return &blogpb.ReadBlogResponse{Blog: &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Title:    data.Title,
		Content:  data.Content,
	}}, nil
}
