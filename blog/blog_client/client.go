package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/blog/blogpb"
	"io"
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
		log.Fatalf("\nUnexpected error: [ %v ]", err)
	}

	fmt.Printf("\nBlog has bem created: [ %v ]", blog.GetBlog())

	readBlog, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blog.GetBlog().GetId()})
	if err != nil {
		log.Fatalf("\nFailed to read blog: [ %v ]", err)
	}

	log.Printf("\nRead blog: [ %v ] ", readBlog)

	blogToUpdate := &blogpb.UpdateBlogRequest{Blog: &blogpb.Blog{
		Id:       readBlog.GetBlog().GetId(),
		AuthorId: "George Updado pra deletar",
		Title:    "Titulo Updado",
		Content:  readBlog.GetBlog().GetContent(),
	}}

	updateBlog, err := c.UpdateBlog(context.Background(), blogToUpdate)
	if err != nil {
		log.Fatalf("Failed to update: [ %v ]", err)
	}

	fmt.Printf("\n Updated blog: [ %v ]", updateBlog.GetBlog().String())

	_, err = c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogToUpdate.GetBlog().GetId()})
	if err != nil {
		log.Fatalf("\nError deleting: [ %v ]", err)
	}

	fmt.Printf("\nDeleted blog: [ %v ]", updateBlog.GetBlog().GetId())

	//	List Blogs

	listBlog, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("\nError calling ListBlog method: [ %v ]", err)
	}

	i := 0
	for {
		res, err := listBlog.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error inside stream for: [ %v ]", err)
		}
		i++
		fmt.Printf("\n[ %v ] -> [ %v ]", i, res.GetBlog().String())
	}
}
