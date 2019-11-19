package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/chengweilun/grpc-go-course/max/maxpb"
	"google.golang.org/grpc"
)

func callMax(c maxpb.MaxClient) {
	requests := []*maxpb.MaxRequest{
		&maxpb.MaxRequest{
			Num: 1,
		},
		&maxpb.MaxRequest{
			Num: 5,
		},
		&maxpb.MaxRequest{
			Num: 3,
		},
		&maxpb.MaxRequest{
			Num: 6,
		},
		&maxpb.MaxRequest{
			Num: 2,
		},
		&maxpb.MaxRequest{
			Num: 20,
		},
	}

	stream, err := c.GetMax(context.Background())
	if err != nil {
		log.Fatalf("can not call rpc GetMax\n")
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			stream.Send(req)
			fmt.Printf("sending request %v\n", req)
			time.Sleep(500 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("error: %v", err)
				close(waitc)
				break
			}
			fmt.Printf("receive message: %v\n", res.GetMaxNum())
		}
	}()

	<-waitc
}

func main() {
	fmt.Println("client")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := maxpb.NewMaxClient(cc)
	fmt.Printf("create client: %v\n", c)
	callMax(c)
}
