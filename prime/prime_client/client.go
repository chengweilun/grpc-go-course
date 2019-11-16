package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/chengweilun/grpc-go-course/prime/primepb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := primepb.NewPrimeServiceClient(cc)
	fmt.Printf("create client: %v\n", c)

	req := &primepb.PrimeRequest{Num: 120}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("404: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while decoding response %v", err)
		}
		fmt.Println(msg.GetPrime())
	}

}
