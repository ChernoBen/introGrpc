package main

import (
	"context"
	"fmt"
	"intro/greet/greetpb"
	"log"

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
	log.Fatalf("Unary Response: %v", response.Result)
}
