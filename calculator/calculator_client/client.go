package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error dialing: [ %v ]", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	//if err := doClientStreaming(c); err != nil {
	//	log.Fatalf("Error calling doClientStreaming: [ %v ]", err)
	//}
	doBiDirectionalStreaming(c)

}

func doBiDirectionalStreaming(c calculatorpb.CalculatorServiceClient) {
	s, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error calling FindMaximum from Server: [ %v ]", err)
		return
	}
	waitc := make(chan struct{})

	sI := []int32{
		1,
		5,
		3,
		6,
		2,
		20,
	}
	go func() {
		fmt.Println("Starting [ send GoRoutine ]")
		//	Sending stream:
		for _, n := range sI {
			if err := s.Send(&calculatorpb.FindMaximumRequest{Number: n}); err != nil {
				close(waitc)
				log.Fatalf("error sending int to server: [ %v ]", err)
				return
			}
		}
		if err := s.CloseSend(); err != nil {
			close(waitc)
			log.Fatalf("error closeSend: [ %v ]", err)
			return
		}
		close(waitc)
	}()

	var rIH []int32
	i := 0
	go func() {
		//	Recieving stream
		for {
			rI, err := s.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				close(waitc)
				log.Fatalf("error recieving data: [ %v ]", err)
				return
			}
			fmt.Printf("Recieved Maximum int: [ %v ]", rI.GetCurrentMaximum())
			rIH[i] = rI.GetCurrentMaximum()
		}
	}()
	<-waitc

	fmt.Printf("All maximums were: [ %v ]", rIH)
}

//func doClientStreaming(c calculatorpb.CalculatorServiceClient) (err error) {
//	s, err := c.ComputeAverage(context.Background())
//	if err != nil {
//		log.Fatalf("Couldn't call ComputeAverage: [ %v ]", err)
//		return
//	}
//	iS := []int64{
//		1,
//		22,
//		33,
//		45,
//	}
//
//	for _, i := range iS {
//		if err = s.Send(&calculatorpb.ComputeRequest{Number: i}); err != nil {
//			log.Fatalf("error sending: [ %v ]", err)
//			return err
//		}
//		log.Printf("Sending via stream: [ %v ] \n", i)
//	}
//	r, err := s.CloseAndRecv()
//	if err != nil {
//		log.Fatalf("error closing and recieving: [ %v ]", err)
//		return err
//	}
//	log.Printf("Recievied as a response an average of: [ %v ]", r.GetAverage())
//	return
//}
