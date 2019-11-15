package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chengweilun/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	fmt.Printf("create client: %v\n", c)

	req := &calculatorpb.CalculatorRequest{
		Add1: 3,
		Add2: 10,
	}

	rep, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("could get resposne: %v", err)
	}

	fmt.Println(rep)
}
