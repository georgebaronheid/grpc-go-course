package main

import (
	"fmt"
    "google.golang.org/grpc"
    "grpc-go-course/prime/primepb"
    "log"
    "net"
    "time"
)

type server struct {
	primepb.UnimplementedPrimeServiceServer
}

func (*server) Decompose(req *primepb.PrimeRequest, stream primepb.PrimeService_DecomposeServer) error {
    input := req.GetInput()
    log.Printf("got input: [ %v ]", input)
    var i int64 = 2
    for input > 1 {
        if input % i == 0 {
            res := &primepb.PrimeResponse{Decomposition: i}
            log.Printf("streaming: [ %v ]", i)
            if err := stream.Send(res); err != nil {
                log.Fatalf("couldn't stream message: [ %v ]", err)
                return err
            }
            input = input / i
        } else {
            i++
        }
        time.Sleep(1000 * time.Millisecond)
    }
	return nil
}

func main() {
	fmt.Println("[ primes ] server up!")

	l, err := net.Listen("tcp", "0.0.0.0:50051")
    if err != nil {
        log.Fatalf("couldnt listen: [ %v ]", err)
    }

    s := grpc.NewServer()
    primepb.RegisterPrimeServiceServer(s, &server{})

    if err := s.Serve(l); err != nil {
        log.Fatalf("couldnt serve: %v", err)
    }
}
