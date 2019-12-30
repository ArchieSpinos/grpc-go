package main

import (
	"context"
	"log"
	"time"

	"github.com/ArchieSpinos/grpc-go/calculator/calcpb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := calcpb.NewCalcServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cc := &calcpb.CalcRequest{
		Numbers: &calcpb.Numbers{
			Num1: 5,
			Num2: 5,
		},
	}

	r, err := c.SumNums(ctx, cc)
	if err != nil {
		log.Fatalf("could not calc: %v", err)
	}

	log.Printf("Calculate sum: %d", r.GetResult())
}
