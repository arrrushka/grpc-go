package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"com.grpc.tleu/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	doGreetingEveryone(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	ctx := context.Background()

	request := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{
		FirstName: "Ilon",
		LastName:  "Musk",
	}}

	response, err := c.Greet(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC %v", err)
	}
	log.Printf("response from Greet:%v", response.Result)
}

func doManyTimesFromServer(c greetpb.GreetServiceClient) {
	ctx := context.Background()
	req := &greetpb.GreetManyTimesRequest{Greeting: &greetpb.Greeting{
		FirstName: "Bob",
		LastName:  "123",
	}}

	stream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we've reached the end of the stream
				break LOOP
			}
			log.Fatalf("error while reciving from GreetManyTimes RPC %v", err)
		}
		log.Printf("response from GreetManyTimes:%v \n", res.GetResult())
	}

}

func doLongGreet(c greetpb.GreetServiceClient) {

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tleu",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Alice",
			},
		},
	}

	ctx := context.Background()
	stream, err := c.LongGreet(ctx)
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doGreetingEveryone(c greetpb.GreetServiceClient) {

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while open stream: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tleu",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Alice",
			},
		},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			log.Printf("Sending message: %v", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending req to server: %v", err.Error())
			}
			time.Sleep(time.Second)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("errow while closing client's stream")
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while getting message from stream: %v", err.Error())
			}
			log.Printf("Received: %v", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}
