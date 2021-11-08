package main

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "grpc-go-course/prime/primepb"
    "io"
    "log"
)

func main() {
    cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("could dial: [ %v ]", err)
    }

    defer cc.Close()

    c := primepb.NewPrimeServiceClient(cc)

    doServerStreaming(c)
}

func doServerStreaming(c primepb.PrimeServiceClient) {
    fmt.Println("[ prime client ] Starting stream request")
    req := &primepb.PrimeRequest{Input: 680}
    fmt.Printf("stream request for %v \n", req)


    stream, err := c.Decompose(context.Background(), req)
    if err != nil {
        log.Fatalf("error requesting stream: [ %v ]", err)
    }
    fmt.Println("starting decomposition:")
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            fmt.Println("End of stream")
            break
        }
        if err != nil {
            log.Fatalf("error receivieng stream? [ %v ]", err)
        }
        mR := msg.GetDecomposition()
        log.Printf("Value: [ %v ]", mR)
    }
}