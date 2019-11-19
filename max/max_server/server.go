package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/chengweilun/grpc-go-course/max/maxpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) GetMax(stream maxpb.Max_GetMaxServer) error {
	fmt.Println("receive max rpc call")
	var maxValue int64 = math.MinInt64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("can't read request\n")
		}
		if req.GetNum() > maxValue {
			senderr := stream.Send(&maxpb.MaxResponse{
				MaxNum: req.GetNum(),
			})
			maxValue = req.GetNum()
			if senderr != nil {
				log.Fatalf("error sending reply  %v\n", err)
				return err
			}
		}
	}
}

func main() {
	fmt.Println("Hello Wolrd!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	maxpb.RegisterMaxServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
