package main

import (
	"context"
	"fmt"
	"intro/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
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

//metodo unary com deadline
func (*server) GreetWithDeadLine(ctx context.Context, req *greetpb.GreetWithDeadLineRequest) (*greetpb.GreetWithDeadLineResponse, error) {
	fmt.Println("greet dealine function invocada!!!")
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("Cliente cancelou a request")
			return nil, status.Error(codes.Canceled, "O cliente cancelou a request!")
		}
		time.Sleep(1 * time.Second)
	}
	//obter dados da request
	first_name := req.GetGreeting().GetFirstName()
	last_name := req.GetGreeting().GetLastName()
	//criando a response
	result := "Hello " + first_name + " " + last_name
	res := &greetpb.GreetWithDeadLineResponse{
		Result: result,
	}
	return res, nil
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

//long greet
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("Init LongGreeting")
	result := ""
	//rodar até o EOF
	for {
		req, err := stream.Recv()
		if err == io.EOF { //fim da linha
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Erro acorrido : %v\n", err)
		}
		first_name := req.GetGreeting().GetFirstName()
		result += "Hello " + first_name + " !"
	}
}

// bi directional method
func (*server) GreetEveryOne(stream greetpb.GreetService_GreetEveryOneServer) error {
	fmt.Println("Bi directional streaming method")
	for {
		req, err := stream.Recv()
		if err != nil {
			log.Fatalf("falha ao ler mensagem: %v\n", err)
			break
		}
		if err == io.EOF {
			break
		}
		first_name := req.GetGreeting().GetFirstName()
		erro := stream.Send(&greetpb.GreetEveryOneResponse{
			Result: "Hello" + first_name,
		})
		if erro != nil {
			log.Fatalf("erro enquanto tentava enviar msg %v\n", erro)
			break
		}
	}
	return nil
}
