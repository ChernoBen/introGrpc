package main

import (
	"context"
	"fmt"
	"intro/greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("greet function invocada!!!")
	//obter dados da request
	first_name := req.GetGreeting().GetFirstName()
	last_name := req.GetGreeting().GetLastName()
	//criando a response
	result := "Hello " + first_name + " " + last_name
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}
func main() {
	fmt.Println("Hello gRPC!!")
	//criar um listener
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Faild to listen: %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to set the server: %v", err)
	}
}
