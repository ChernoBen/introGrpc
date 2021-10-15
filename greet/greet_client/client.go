package main

import (
	"context"
	"fmt"
	"intro/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello client!!")
	//1# - criar uma conexão com servidor
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect : %v", err)
	}
	defer conn.Close()
	//2# - criar um client
	c := greetpb.NewGreetServiceClient(conn)
	if c != nil {
		fmt.Println("Client created")
	}
	doUnary(c)
	doServerStreaming(c)
	doClientStreaming(c)
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
