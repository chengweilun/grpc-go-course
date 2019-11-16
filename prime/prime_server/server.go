package main

import (
	"fmt"
	"log"
	"net"

	"github.com/chengweilun/grpc-go-course/prime/primepb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) PrimeNumberDecomposition(req *primepb.PrimeRequest, stream primepb.PrimeService_PrimeNumberDecompositionServer) error {
	num := req.GetNum()
	var k int64 = 2
	for num > 1 {
		if num%k == 0 {
			stream.Send(&primepb.PrimeResponse{Prime: k})
			num = num / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func main() {
	fmt.Println("Hello Wolrd!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	primepb.RegisterPrimeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
