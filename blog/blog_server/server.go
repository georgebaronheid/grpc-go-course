package main

import (
    "fmt"
    "google.golang.org/grpc"
    "grpc-go-course/blog/blogpb"
    "log"
    "net"
    "os"
    "os/signal"
)

type server struct {
    blogpb.UnimplementedBlogServiceServer
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
    fmt.Println("Graceful shutdown done")

}
