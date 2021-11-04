package main

import (
	"context"
	"google.golang.org/grpc"
	"grpc-go-course/sum/sumpb"
	"log"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[ client ] coudn't create conn: %v", err)
	}

	defer cc.Close()

	c := sumpb.NewSumServiceClient(cc)

	req := &sumpb.SumRequest{FirstNumber: 1, SecondNumber: 123}
	log.Println(req.String())

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("[ client ] couldn't get Sum: %v", err)
	}

	log.Println("Sum result: " + res.String())

}
