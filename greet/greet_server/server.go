package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/greet/greetpb"
	"log"
	"net"
	"time"
)

type server struct {
	//Required func
	greetpb.UnimplementedGreetServiceServer
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fN, lN := req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName()
	for i := 1; i <= 10; i++ {
		r := "Hello " + fN + " " + lN + ". Number [ " + string(rune(i)) + "] "
		res := &greetpb.GreetManyTimesResponse{Result: r}
		err := stream.Send(res)
		if err != nil {
			log.Fatalf("Couldn't send message: [ %v ]", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)
	f, l := req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName()
	result := "Hello " + f + " " + l
	res := greetpb.GreetResponse{Result: result}
	return &res, nil
}

func main() {
	fmt.Println("[ greet_server ] Up!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}
