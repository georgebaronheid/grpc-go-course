package main

import (
	"context"
    "fmt"
    "google.golang.org/grpc"
	"grpc-go-course/greet/greetpb"
	"log"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
    doUnaryCall(c)
}

func doUnaryCall(c greetpb.GreetServiceClient) {
    fmt.Println("Starting UNARY RPC Call")
    req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "George", LastName: "Baronheid"}}
    res, err := c.Greet(context.Background(), req)
    if err != nil {
        log.Fatalf("Error while calling greet RPC: %v", err)
    }
    log.Printf("Response from greet RPC: %v", res.String())
}
