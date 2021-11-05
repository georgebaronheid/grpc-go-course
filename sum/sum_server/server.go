package main

import (
    "context"
    "google.golang.org/grpc"
    "grpc-go-course/sum/sumpb"
    "log"
    "net"
)

type server struct {
    sumpb.UnimplementedSumServiceServer
}

func (*server) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error){
    n1, n2 := req.GetFirstNumber(), req.GetSecondNumber()
    s := n1 + n2
    pbR := &sumpb.SumResponse{Result: s}
    return pbR, nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
    if err != nil {
        log.Fatalf("[ server ] listener error: %v", err)
    }

    s := grpc.NewServer()
    sumpb.RegisterSumServiceServer(s, &server{})

    if err := s.Serve(l); err != nil {
        log.Fatalf("[ server ] couldn't serve: %v", err)
    }

}
