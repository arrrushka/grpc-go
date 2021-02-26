package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"com.grpc.tleu/calculator/calcpb"
	"google.golang.org/grpc"
)

//Server with embedded UnimplementedCalculatorServiceServer
type Server struct {
	calcpb.UnimplementedCalculatorServiceServer
}

// PrimeNumberDecomposition is an example of stream from server side
func (s *Server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest,
	stream calcpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v \n", req)
	number := req.GetNumber()
	for number > 1 {
		num := getFirstPrime(number)
		number /= num
		res := &calcpb.PrimeNumberDecompositionResponse{Number: num}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error while sending greet many times responses: %v", err.Error())
		}
		time.Sleep(time.Second)
	}
	return nil
}

func getFirstPrime(number int32) int32 {
	for i := 2; int32(i) <= number; i++ {
		if number%int32(i) == 0 {
			return int32(i)
		}
	}
	return number
}

func (s *Server) ComputeAverage(stream calcpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Average function was invoked with a streaming request\n")
	var result int32
	result = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Number: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		number := req.GetNumber()
		result += number
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	calcpb.RegisterCalculatorServiceServer(s, &Server{})
	log.Println("Server is running on port:50052")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
