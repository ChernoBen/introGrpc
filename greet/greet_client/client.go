package main

import (
	"context"
	"fmt"
	"intro/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello client!!")
	//passando credenciais para comunicação com servidor
	certFile := "ssl/ca.crt"
	creds, err := credentials.NewServerTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("Erro ao obter credenciais: %v\n", err)
	}
	opts := grpc.WithTransportCredentials(creds)
	//1# - criar uma conexão com servidor
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to connect : %v", err)
	}
	defer conn.Close()
	//2# - criar um client
	c := greetpb.NewGreetServiceClient(conn)
	if c != nil {
		fmt.Println("Client created")
	}
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doUnaryWithDeadLine(c, 5*time.Second)
	doUnaryWithDeadLine(c, 1*time.Second)
}

//unary com deadline
func doUnaryWithDeadLine(c greetpb.GreetServiceClient, seconds time.Duration) {
	fmt.Println("Inicio func unaria")
	//criando request
	req := &greetpb.GreetWithDeadLineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Benjamim",
			LastName:  "Francisco",
		},
	}
	//definindo deadline no context
	ctx, cancel := context.WithTimeout(context.Background(), seconds)
	defer cancel()
	//invocando metodo unario
	response, error := c.GreetWithDeadLine(ctx, req)
	if error != nil {
		statusErr, ok := status.FromError(error)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Tempo maximo atingido")
			} else {
				fmt.Printf("Erro inesperado: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Error enquanto chamava a Greet %v \n", error)
		}
		return
	}
	log.Printf("Unary Response: %v", response.Result)
}

// func que consulta BiDiStreaming
func doBiDiStreaming(c greetpb.GreetServiceClient) {
	//criar um streaming invocando o client
	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalf("Error enquanto esperava resposta do server %v\n", err)
		return
	}
	// criando requests
	requests := []*greetpb.GreetEveryOneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Ben",
				LastName:  "ja",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Francisco",
				LastName:  "Benjamim",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Cherno",
				LastName:  "Ben",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Benjamim",
				LastName:  "Francisco",
			},
		},
	}
	//criando chanel de espera
	waitch := make(chan struct{})
	//usar go routine para receber as mensagens do servidor
	go func() {
		//to send
		for _, req := range requests {
			fmt.Printf("Mensagem enviada: %v\n", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("erro ao enviar msg %v\n", err)
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		//parar envio
		stream.CloseSend()
	}()
	go func() {
		//to receive
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				//se EOF break
				break
			}
			if err != nil {
				log.Fatalf("Falha ao receber mensagem: %v\n", err)
				break
			}
			fmt.Printf("Respostas : %v\n", res.GetResult())
		}
		close(waitch)
	}()
	//manter bloqueado até que as mensagens terminarem
	<-waitch
}

//func que chama metodo client streaming
func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Client streaming RPC")
	//criando as requests
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Ben",
				LastName:  "ja",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Francisco",
				LastName:  "Benjamim",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Cherno",
				LastName:  "Ben",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Benjamim",
				LastName:  "Francisco",
			},
		},
	}
	//nesse tipo, client streaming não é necessario passar request, apenas context.Background
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Erro na chamada da funcao client streaming: %v\n", err)
	}
	//iterar o slice de requests
	for _, req := range requests {
		fmt.Printf("request sent : %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	//fechando stream
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("erro ao receber msg de LongGreeting %v\n", err)
	}
	fmt.Printf("Response: %v", res)
}

//funcao que chama metodo unario
func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Inicio func unaria")
	//criando request
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Benjamim",
			LastName:  "Francisco",
		},
	}
	//invocando metodo unario
	response, error := c.Greet(context.Background(), req)
	if error != nil {
		log.Fatalf("Error enquanto chamava a Greet %v \n", error)
	}
	log.Printf("Unary Response: %v", response.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Server streaming Rpc...")
	//criando greetManyTimes request
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Benja",
			LastName:  "Francisco",
		},
	}
	resStream, error := c.GreetManyTimes(context.Background(), req)
	if error != nil {
		log.Fatalf("Deus Ruim %v \n", error)
	}
	for {
		msg, err := resStream.Recv()
		//se errro igual a endOfFile : break
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Falha na leitura da stream %v \n", error)
		}
		log.Printf("Resposta de GreetManyTime: %v\n", msg.GetResult())
	}
}
