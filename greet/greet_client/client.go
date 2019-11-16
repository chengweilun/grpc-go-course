package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/chengweilun/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func unaryRequest(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "jane", LastName: "doe",
		},
	}

	rep, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("could get resposne: %v", err)
	}

	fmt.Println(rep)

}

func serverStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "jane", LastName: "doe",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while server streaming %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// reach end
			break
		}
		if err != nil {
			log.Fatalf("err while read response %v", err)
		}
		fmt.Println(msg.GetResult())
	}

}

func main() {
	fmt.Println("client")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("create client: %v\n", c)
	unaryRequest(c)
	serverStreaming(c)
}
