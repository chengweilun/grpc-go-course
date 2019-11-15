package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/chengweilun/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorReponse, error) {
	add1 := req.GetAdd1()
	add2 := req.GetAdd2()
	result := add1 + add2
	return &calculatorpb.CalculatorReponse{Result: result}, nil
}

func main() {
	fmt.Println("Hello Wolrd!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
