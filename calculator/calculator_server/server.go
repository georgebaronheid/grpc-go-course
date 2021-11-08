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
		s, err := stream.Recv()
		tR += s.GetNumber()
		if err == io.EOF {
			log.Printf("Total recieved: [ %v ], ammount recieved: [ %v ]", tR, k)
			a = float64(tR) / k
			if err := stream.SendAndClose(&calculatorpb.ComputeResponse{Average: a}); err != nil {
				log.Fatalf("error send and close: [ %v ]", err)
			}
            return nil
		} else if err != nil {
			log.Fatalf("error recieving stream: [ %v ]", err)
			return err
		}
		log.Printf("Recieved: [ %v ] \n", s)
		k++
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
