package main

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "grpc-go-course/calculator/calculatorpb"
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
	//doBiDirectionalStreaming(c)

	doSqrtWithError(c)
}

func doSqrtWithError(cSC calculatorpb.CalculatorServiceClient) {
	var c1, c2 *calculatorpb.SquareRootResponse
	var err error
	//Correct calls
	if c1, err = cSC.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: 9.0}); err != nil {
        errorHandling(err)
    }
    fmt.Printf("Squareroot of 9.0: [ %v ]\n", c1.GetNumberRoot())

    if c2, err = cSC.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: 0.0}); err != nil {
        errorHandling(err)
    }
    fmt.Printf("Squareroot of 0.0: [ %v ]\n", c2.GetNumberRoot())


    //Incorrect call
    if _, err = cSC.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: -9.0}); err != nil {
        errorHandling(err)
    }
}

func errorHandling(err error) {
    respErr, ok := status.FromError(err)
    if ok {
        //    Known error
        fmt.Println(respErr.Code())
        fmt.Println(respErr.Message())
        if respErr.Code() == codes.InvalidArgument {
            fmt.Println("We probably sent an incorrect number")
        }
    } else {
        log.Fatalf("Unknow error [ %v ]", err)
    }
}

//func doBiDirectionalStreaming(cSC calculatorpb.CalculatorServiceClient) {
//	fMC, err := cSC.FindMaximum(context.Background())
//	if err != nil {
//		log.Fatalf("error calling FindMaximum from Server: [ %v ]", err)
//		return
//	}
//	waitc := make(chan struct{})
//
//	iS := []int32{
//		1,
//		5,
//		3,
//		6,
//		2,
//		20,
//		200,
//		200,
//		2000,
//		2000,
//		2001,
//	}
//	go func() {
//		fmt.Println("Starting [ send GoRoutine ]")
//		//	Sending stream:
//		for _, n := range iS {
//			if err := fMC.Send(&calculatorpb.FindMaximumRequest{Number: n}); err != nil {
//				close(waitc)
//				log.Fatalf("error sending int to server: [ %v ]", err)
//			}
//		}
//		if err := fMC.CloseSend(); err != nil {
//			close(waitc)
//			log.Fatalf("error closeSend: [ %v ]", err)
//		}
//	}()
//
//	var rIH []int32
//	go func() {
//		//	Receiving stream
//		for {
//			rI, err := fMC.Recv()
//			if err == io.EOF {
//				break
//			}
//			if err != nil {
//				close(waitc)
//				log.Fatalf("error recieving data: [ %v ]", err)
//			}
//			fmt.Printf("Recieved Maximum int: [ %v ]\n", rI.GetCurrentMaximum())
//			rIH = append(rIH, rI.GetCurrentMaximum())
//		}
//		close(waitc)
//	}()
//	<-waitc
//
//	fmt.Printf("All maximums were: %v\n", rIH)
//}

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
