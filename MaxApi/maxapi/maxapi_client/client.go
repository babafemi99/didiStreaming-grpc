package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"maxapi/maxapi/maxapipb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func doBiDiStreaming(c maxapipb.GetMaxApiClient) {
	fmt.Println("Starting to bi directional streaming...")
	stream, err := c.MaxApi(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}
	waitc := make(chan struct{})

	go func() {
		numbers := []int64{1, 6, 8, 3, 9, 11, 4, 20}
		for _, num := range numbers {
			fmt.Printf("sending message:%v \n", num)
			stream.Send(&maxapipb.MaxApiRequest{
				Number: num,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		defer close(waitc) 
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving stream: %v", err)
				break
			}
			fmt.Printf("Recieved, New Max is: %v\n", res.GetMaxNumber())
		}
	}()
	<-waitc
}
func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	// c := bidipb.NewGreetServiceClient(cc)
	c := maxapipb.NewGetMaxApiClient(cc)
	doBiDiStreaming(c)
}
