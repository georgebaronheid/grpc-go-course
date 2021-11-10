package main

import (
	"context"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"log"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error dialing: [ %v ]", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	if err := doClientStreaming(c); err != nil {
		log.Fatalf("Error calling doClientStreaming: [ %v ]", err)
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) (err error) {
	s, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Couldn't call ComputeAverage: [ %v ]", err)
		return
	}
	iS := []int64{
		1,
		22,
		33,
		45,
	}

	for _, i := range iS {
		if err = s.Send(&calculatorpb.ComputeRequest{Number: i}); err != nil {
			log.Fatalf("error sending: [ %v ]", err)
			return err
		}
		log.Printf("Sending via stream: [ %v ] \n", i)
	}
	r, err := s.CloseAndRecv()
	if err != nil {
		log.Fatalf("error closing and recieving: [ %v ]", err)
		return err
	}
	log.Printf("Recievied as a response an average of: [ %v ]", r.GetAverage())
	return
}
