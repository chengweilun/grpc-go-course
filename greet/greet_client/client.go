package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

func clientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane",
				LastName:  "don",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane2",
				LastName:  "don",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane3",
				LastName:  "don",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane4",
				LastName:  "don",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane5",
				LastName:  "don",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane6",
				LastName:  "don",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("can not call rpc longGreet\n")
	}

	for _, req := range requests {
		fmt.Printf("sending request %v\n", req)
		stream.Send(req)
		time.Sleep(500 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("can not get response from server")
	}
	fmt.Printf("%v\n", res.GetResult())
}

func biStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane",
				LastName:  "don",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane2",
				LastName:  "don",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane3",
				LastName:  "don",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane4",
				LastName:  "don",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane5",
				LastName:  "don",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jane6",
				LastName:  "don",
			},
		},
	}
	waitc := make(chan struct{})
	stream, err := c.GreetEveryOne(context.Background())

	if err != nil {
		log.Fatalf("can not call rpc method GreetEveryOne")
	}

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
			fmt.Printf("receive message: %v\n", res.GetResult())
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

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("create client: %v\n", c)
	biStreaming(c)
}
