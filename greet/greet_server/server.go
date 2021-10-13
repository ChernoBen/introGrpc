package main

import (
	"context"
	"fmt"
	"intro/greet/greetpb"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

//metodo unary para struct server
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

//metodo server streaming p/ struct server
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("greetManyTimes foi invocada com sucesso")
	firstName := req.GetGreeting().GetFirstName()
	//criando laço que responderá x vezes
	for i := 0; i < 10; i++ {
		result := "Olá " + firstName + " numero " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		//enviando N vezes a response
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
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
