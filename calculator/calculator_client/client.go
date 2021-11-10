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

	defer func(cc *grpc.ClientConn) {
		if err := cc.Close(); err != nil {
			log.Fatalf("error closing ClientConn: [%v]", err)
		}
	}(cc)

	c := calculatorpb.NewCalculatorServiceClient(cc)

	//if err := doClientStreaming(c); err != nil {
	//	log.Fatalf("Error calling doClientStreaming: [ %v ]", err)
	//}
	doBiDirectionalStreaming(c)

}

func doBiDirectionalStreaming(cSC calculatorpb.CalculatorServiceClient) {
	fMC, err := cSC.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error calling FindMaximum from Server: [ %v ]", err)
		return
	}
	waitc := make(chan struct{})

	iS := []int32{
		1,
		5,
		3,
		6,
		2,
		20,
		200,
		200,
		2000,
		2000,
		2001,
	}
	go func() {
		fmt.Println("Starting [ send GoRoutine ]")
		//	Sending stream:
		for _, n := range iS {
			if err := fMC.Send(&calculatorpb.FindMaximumRequest{Number: n}); err != nil {
				close(waitc)
				log.Fatalf("error sending int to server: [ %v ]", err)
			}
		}
		if err := fMC.CloseSend(); err != nil {
			close(waitc)
			log.Fatalf("error closeSend: [ %v ]", err)
		}
	}()

	var rIH []int32
	go func() {
		//	Receiving stream
		for {
			rI, err := fMC.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				close(waitc)
				log.Fatalf("error recieving data: [ %v ]", err)
			}
			fmt.Printf("Recieved Maximum int: [ %v ]\n", rI.GetCurrentMaximum())
			rIH = append(rIH, rI.GetCurrentMaximum())
		}
		close(waitc)
	}()
	<-waitc

	fmt.Printf("All maximums were: %v\n", rIH)
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
