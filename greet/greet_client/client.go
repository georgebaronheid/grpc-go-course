package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/greet/greetpb"
    "io"
    "log"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//doUnaryCall(c)

	doServerStreaming(c)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("[ greet client ] Starting stream request")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
            FirstName: "George",
            LastName: "Baronheid",
        },
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
    if err != nil {
        log.Fatalf("[ client ] error while streaming GreetManyTimes: %v", err)
    }
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            //End of stream
            break;
        }
        if err != nil {
            log.Fatalf("Error while reading stream: %v", err)
        }
        log.Printf("Results from GreetManyTimes: %v", msg.GetResult())
    }

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
