package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"com.grpc.tleu/calculator/calcpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := calcdb.NewCalculatorServiceClient(conn)
	doAverage(c)
}

func doAverage(c calcdb.CalculatorServiceClient) {

	requests := []*calcdb.ComputeAverageRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
	}

	ctx := context.Background()
	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}

func doManyTimesFromServer(c calcdb.CalculatorServiceClient) {
	ctx := context.Background()
	req := &calcdb.PrimeNumberDecompositionRequest{Number: 120}

	stream, err := c.PrimeNumberDecomposition(ctx, req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC %v", err)
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
			log.Fatalf("error while reciving from PrimeNumberDecomposition RPC %v", err)
		}
		log.Printf("response from PrimeNumberDecomposition:%v \n", res.GetNumber())
	}

}
