package main

import (
	"fmt"
	"io"
	"log"
	"maxapi/maxapi/maxapipb"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	maxapipb.UnimplementedGetMaxApiServer
}

func (*server) MaxApi(stream maxapipb.GetMaxApi_MaxApiServer) error {
	fmt.Println("Recieved MaxApi")
	maximum := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		number := req.GetNumber()
		if number > maximum {
			maximum = number
			sendErr := stream.Send(&maxapipb.MaxApiResponse{
				MaxNumber: maximum,
			})
			if sendErr != nil {
				log.Fatalf("Error while reading client stream: %v", sendErr)
				return sendErr
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	maxapipb.RegisterGetMaxApiServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
