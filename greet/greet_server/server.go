package main

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/status"
    "grpc-go-course/greet/greetpb"
    "io"
    "log"
    "net"
    "time"
)

type server struct {
    //Required func
    greetpb.UnimplementedGreetServiceServer
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
    fN, lN := req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName()
    for i := 1; i <= 10; i++ {
        r := "Hello " + fN + " " + lN + ". Number [ " + string(rune(i)) + "] "
        res := &greetpb.GreetManyTimesResponse{Result: r}
        err := stream.Send(res)
        if err != nil {
            log.Fatalf("Couldn't send message: [ %v ]", err)
        }
        time.Sleep(1000 * time.Millisecond)
    }
    return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
    r := "Hello"
    for {
        s, err := stream.Recv()
        if err == io.EOF {
            if err := stream.SendAndClose(&greetpb.LongGreetResponse{Result: r}); err != nil {
                log.Fatalf("error sending and closing stream: [ %v ]", err)
                return err
            }
            return nil
        }
        if err != nil {
            log.Fatalf("error while reading client stream: [ %v ]", err)
            return err
        }
        r += s.GetGreeting().FirstName + "! "
    }
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
    fmt.Printf("Greet function was invoked with %v", req)
    f, l := req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName()
    result := "Hello " + f + " " + l
    res := greetpb.GreetResponse{Result: result}
    return &res, nil
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
    fmt.Printf("Greet function was invoked with %v", req)
    for i := 0; i < 3; i++ {
        if ctx.Err() == context.Canceled {
            fmt.Println("Client side cancel")
            return nil, status.Errorf(codes.Canceled, "Client side request cancel")
        }
        time.Sleep(1 * time.Second)
    }
    f, l := req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName()
    result := "Hello " + f + " " + l
    res := greetpb.GreetWithDeadlineResponse{Result: result}
    return &res, nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
    fmt.Printf("GreetEveryone function was invoked")
    const greeting = "Hello "
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            log.Fatalf("error recieving stream [ %v ]", err)
            return err
        }
        r := &greetpb.GreetEveryoneResponse{Result: greeting + req.GetGreeting().GetFirstName()}
        if err := stream.Send(r); err != nil {
            log.Fatalf("error sending stream response [ %v ]", err)
            return err
        }
        time.Sleep(500 * time.Millisecond)
    }
}

func main() {
    fmt.Println("[ greet_server ] Up!")

    certFile := "ssl/server.crt"
    keyFile := "ssl/server.pem"
    creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
    if err != nil {
        log.Fatalf("Failed loading certificates [ %v ]", err)
        return
    }

    opts := grpc.Creds(creds)

    s := grpc.NewServer(opts)

    lis, err := net.Listen("tcp", "0.0.0.0:50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    greetpb.RegisterGreetServiceServer(s, &server{})

    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to server: %v", err)
    }

}
