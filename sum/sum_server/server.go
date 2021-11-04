package main

import (
    "context"
    "google.golang.org/grpc"
    "grpc-go-course/sum/sumpb"
    "log"
    "net"
)

type serverSum struct {
    sumpb.UnimplementedSumServiceServer
}

func (*serverSum) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error){
    v1, v2 := req.GetFirstNumber(), req.GetSecondNumber()
    r := v1 + v2
    grpcR := &sumpb.SumResponse{Result: r}
    return grpcR, nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
    if err != nil {
        log.Fatalf("[ serverSum ] listener error: %v", err)
    }

    s := grpc.NewServer()
    sumpb.RegisterSumServiceServer(s, &serverSum{})

    if err := s.Serve(l); err != nil {
        log.Fatalf("[ serverSum ] couldn't serverSum: %v", err)
    }

}
