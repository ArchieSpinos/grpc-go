package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ArchieSpinos/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number" + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.SendMsg(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) CalcPrimes(req *greetpb.CalcPrimesRequest, streamNumbers greetpb.GreetService_CalcPrimesServer) error {
	fmt.Printf("CalcPrimes function was invoked with initial value %v,", req)
	var (
		k int64 = 2
		N int64 = req.GetCalcprimes().GetNumberToCalc()
	)

	for N > 1 {
		if N%k == 0 { // if k evenly divides into N
			res := &greetpb.CalcPrimesResponse{
				Result: k,
			}
			streamNumbers.SendMsg(res)
			N = N / k // divide N by k so that we have the rest of the number left.
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("calling LongGreet function")
	result := "Hello"
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += firstName
	}
}

func (*server) ComputeAverage(stream greetpb.GreetService_ComputeAverageServer) error {
	fmt.Printf("Calling ComputeAverage function")
	var (
		total float64 = 0
		sum   float64 = 0
	)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.ComputeAverageResponse{
				Avg: float64(total / sum),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum++
		total += req.GetNum()
		fmt.Println(sum, total)
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading bidi stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending bidi stream to client: %v", err)
		}
	}
}

func (*server) FindMaximum(stream greetpb.GreetService_FindMaximumServer) error {
	var max float64 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading bidi stream: %v", err)
		}
		if current := req.GetNum(); current > max {
			stream.Send(&greetpb.FindMaximumResponse{
				Num: current,
			})
			max = current
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *greetpb.SquareRootRequest) (*greetpb.SquareRootResponse, error) {
	fmt.Printf("Calling SquareRoot function")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &greetpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func (*server) WithDeadline(ctx context.Context, req *greetpb.WithDeadlineRequest) (*greetpb.WithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with %v", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.WithDeadlineResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
