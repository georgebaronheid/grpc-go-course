package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/blog/blogpb"
	"log"
)

func main() {
	opts := grpc.WithInsecure()

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Error closing connection inside defer: [ %v ]", err)
		}
	}(conn)

	c := blogpb.NewBlogServiceClient(conn)

	in := blogpb.Blog{
		AuthorId: "George",
		Title:    "Blog do Jorjinho",
		Content:  "Aqui tem um belo de um conteúdo",
	}

	blog, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: &in})
	if err != nil {
		log.Fatalf("Unexpected error: [ %v ]", err)
	}

	fmt.Printf("Blog has bem created: [ %v ]\n\n", blog.GetBlog())

	readBlog, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blog.GetBlog().GetId()})
	if err != nil {
		log.Fatalf("Failed to read blog: [ %v ]", err)
	}

	log.Printf("Read blog: [ %v ] ", readBlog)
}
