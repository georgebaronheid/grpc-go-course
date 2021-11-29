package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage called!")
	tR := int64(0)
	k := 0.0
	var a float64
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Total recieved: [ %v ], ammount recieved: [ %v ]", tR, k)
			a = float64(tR) / k
			if err = stream.SendAndClose(&calculatorpb.ComputeResponse{Average: a}); err != nil {
				log.Fatalf("error send and close: [ %v ]", err)
				return err
			}
			return nil
		} else if err != nil {
			log.Fatalf("error recieving stream: [ %v ]", err)
			return err
		}
		log.Printf("Recieved: [ %v ] \n", r)
		tR += r.GetNumber()
		k++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum was invoked")
	m := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error recieving request: [ %v ]", err)
			return err
		}
		log.Printf("Recieved: [ %v ], current maximum: [ %v ]", req.GetNumber(), m)
		r := req.GetNumber()
		if r > m {
			if err := stream.Send(&calculatorpb.FindMaximumResponse{CurrentMaximum: r}); err != nil {
				log.Fatalf("error sending message: [ %v ]", err)
				return err
			}
			m = r
		}
	}
}

func (*server) SquareRoot(ctx context.Context, in *calculatorpb.SquareRootRequest) (out *calculatorpb.SquareRootResponse, err error) {
	req := in.GetNumber()

	if req < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("[calculator] Received a negative number as an input: [%v]", req))
	}
	if req == 0 {
		out = &calculatorpb.SquareRootResponse{NumberRoot: 0.0}
		return
	}

	out = &calculatorpb.SquareRootResponse{NumberRoot: math.Sqrt(float64(req))}

	return
}

func main() {
	fmt.Println("[ calculator ] Server up!")
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("couldn't listen: [ %v ]", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(l); err != nil {
		log.Fatalf("couldn't serve: [ %v ]", err)
	}
}
