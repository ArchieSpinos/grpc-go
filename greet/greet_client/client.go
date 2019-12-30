package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/ArchieSpinos/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("created client %f", c)
	// doUnary(c)

	// doServerStreaming(c)

	// getPrimes(c)
	//doClientStreaming(c)
	//getAverage(c)
	// doBiDiStreaming(c)
	// findMaxStreaming(c)
	// doErrorUnary(c)
	doUnaryWithDeadline(c, 1*time.Second)
	doUnaryWithDeadline(c, 5*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Archie",
			LastName:  "Spinos",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Server Streaming...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Archie",
			LastName:  "Spinos",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func getPrimes(c greetpb.GreetServiceClient) {
	fmt.Println("Starting prime number decompossion...")

	req := &greetpb.CalcPrimesRequest{
		Calcprimes: &greetpb.CalcPrimes{
			NumberToCalc: 120,
		},
	}
	calcStream, err := c.CalcPrimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error whole calling CalcPrimes: %v", err)
	}

	for {
		msg, err := calcStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("These are the prime numbers decomposed: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Client Streaming...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Archie",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tom",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hogn",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mary",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Thanos",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error calling LongGreet: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while recieving response: %v", err)
	}
	fmt.Printf("LongGreet response: %v\n", res)

}

func getAverage(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do calculate average...")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error calling ComputeAverage: %v", err)
	}

	requests := []float64{2, 4, 6, 8}

	for _, req := range requests {
		fmt.Printf("sending numbers for avg calculation\n")
		stream.Send(&greetpb.ComputeAverageRequest{
			Num: req,
		})
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while recieving response: %v", err)
	}
	fmt.Printf("ComputeAverage response: %v\n", res.Avg)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
	}

	Names := []string{"Archie", "Tom", "John", "Mary", "Jim"}

	waitc := make(chan struct{})

	go func() {
		for _, name := range Names {
			fmt.Printf("sending message to: %v", name)
			stream.Send(&greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: name,
				},
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}

func findMaxStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
	}

	nums := []float64{1, 4, 3, 6, 8, 3, 3, 24, 56, 67, 77, 87, 4, 1, 0}

	waitc := make(chan struct{})

	go func() {
		for _, num := range nums {
			stream.Send(&greetpb.FindMaximumRequest{
				Num: num,
			})
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("This is the new max: %v", res.GetNum())
		}
		close(waitc)
	}()

	<-waitc
}

func doErrorUnary(c greetpb.GreetServiceClient) {
	fmt.Print("Starting to do SquareRoot Unary RPC")

	doErrorCall(c, 9)
	doErrorCall(c, -3)
}

func doErrorCall(c greetpb.GreetServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &greetpb.SquareRootRequest{
		Number: n,
	})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("we probably sent a negative number")
				return
			}
		} else {
			log.Fatalf("Big error calling squareRoot: %v", err)
			return
		}
	}
	fmt.Printf("Square root of %v is: %v\n", n, res.GetNumberRoot())
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting a Unary with deadline RPC...")
	req := &greetpb.WithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Archie",
			LastName:  "Spinos",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.WithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("Unexcpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling WithDeadline RPC: %v", err)
		}
		return
	}
	log.Printf("Response from WithDeadline: %v", res.Result)
}
