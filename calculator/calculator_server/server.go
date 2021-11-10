package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
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
	for  {
		req, err := stream.Recv()
		log.Printf("Recieved: [ %v ], current maximum: [ %v ]", req.GetNumber(), m)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error recieving request: [ %v ]", err)
			return err
		}
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
