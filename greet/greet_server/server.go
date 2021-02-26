package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"com.grpc.tleu/greet/greetpb"
	"google.golang.org/grpc"
)

//Server with embedded UnimplementedGreetServiceServer
type Server struct {
	greetpb.UnimplementedGreetServiceServer
}

//Greet is an example of unary rpc call
func (s *Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v \n", req)
	firstName := req.GetGreeting().GetFirstName()

	result := "Hello, " + firstName

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

// GreetManyTimes is an example of stream from server side
func (s *Server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v \n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{Result: fmt.Sprintf("%d) Hello, %v\n", i, firstName)}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error while sending greet many times responses: %v", err.Error())
		}
		time.Sleep(time.Second)
	}
	return nil
}

//LongGreet is an example of stream from client side
func (s *Server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	var result string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})

		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		firstName := req.Greeting.GetFirstName()
		result += "Hello " + firstName + "! \n"
	}
}

//GreetEveryone is an example of bidirectional stream
func (s *Server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err.Error())
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello, " + firstName
		err = stream.Send(&greetpb.GreetEveryoneResponse{Result: result})
		if err != nil {
			log.Fatalf("error while sending to client: %v", err.Error())
			return err
		}
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
