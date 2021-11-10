package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//doUnaryCall(c)

	//doServerStreaming(c)
	//doClientStreaming(c)
	doClientBiStreaming(c)

}

func doClientBiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Bidi streaming RPC")
	/**
	We create a stream by invoking the client, send a bunch o messages to the server by go routine,
	recieve a bunch of mesages from the client and block until everything is done
	*/

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "George",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Arroz",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Feijão",
			},
		},
	}

	s, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling stream: [ %v ]", err)
		return
	}

	waitc := make(chan struct{})
	go func() {
		fmt.Println("Starting [ Send GoRoutine ]")
		//	function to send a bunch o message
		for _, req := range requests {
			fmt.Printf("Sending message: [ %v ]\n", req)
			if err := s.Send(req); err != nil {
				log.Fatalf("error sending message")
				return
			}
			time.Sleep(250 * time.Millisecond)
		}
		if err := s.CloseSend(); err != nil {
			log.Fatalf("error closesend: [ %v ]", err)
		}
	}()

	go func() {
		fmt.Println("Starting [ Recieve GoRoutine ]")
		//	function to send a bunch o message
		for {
			r, err := s.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while recieving: [ %v ]", err)
				break
			}
			log.Printf("Recieving data: [ %v ]\n", r)
		}
		close(waitc)
	}()

	//block until everything is done
	<-waitc
}

//func doClientStreaming(c greetpb.GreetServiceClient) {
//	s, err := c.LongGreet(context.Background())
//	if err != nil {
//		log.Fatalf("error while calling longgreet [ %v ]", err)
//	}
//	requests := []*greetpb.LongGreetRequest{
//		&greetpb.LongGreetRequest{
//			Greeting: &greetpb.Greeting{
//				FirstName: "George",
//			},
//		},
//		&greetpb.LongGreetRequest{
//			Greeting: &greetpb.Greeting{
//				FirstName: "Arroz",
//			},
//		},
//		&greetpb.LongGreetRequest{
//			Greeting: &greetpb.Greeting{
//				FirstName: "Feijão",
//			},
//		},
//	}
//	for _, request := range requests {
//		fmt.Printf("Sending: [ %v ] \n", request.GetGreeting().GetFirstName())
//		s.Send(request)
//	}
//	r, err := s.CloseAndRecv()
//	if err != nil {
//		log.Fatalf("couldnt close and recieve: [ %v ]", err)
//	}
//	log.Printf("response: %v \n", r.GetResult())
//}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("[ greet client ] Starting stream request")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "George",
			LastName:  "Baronheid",
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
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Results from GreetManyTimes: %v", msg.GetResult())
	}

}

//func doUnaryCall(c greetpb.GreetServiceClient) {
//	fmt.Println("Starting UNARY RPC Call")
//	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "George", LastName: "Baronheid"}}
//	res, err := c.Greet(context.Background(), req)
//	if err != nil {
//		log.Fatalf("Error while calling greet RPC: %v", err)
//	}
//	log.Printf("Response from greet RPC: %v", res.String())
//}
