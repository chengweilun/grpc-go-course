package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/chengweilun/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetReponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := "hello" + firstName
	return &greetpb.GreetReponse{Result: result}, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesReponse{Result: strconv.Itoa(i)}
		stream.Send(res)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// finish read request
			return stream.SendAndClose(&greetpb.LongGreetReponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading request streaming: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "hello " + firstName + "\n"
	}
}

func main() {
	fmt.Println("Hello Wolrd!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
