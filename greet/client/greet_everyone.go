package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/arthsan/grpc-go/greet/proto"
)

func doGreetEveryOne(c pb.GreetServiceClient) {
	log.Println("doGreetEveryOne was invoked")

	stream, err := c.GreetEveryOne(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream %v\n", err)
	}

	reqs := []*pb.GreetRequest{
		{FirstName: "Arthsan"},
		{FirstName: "Daniela"},
		{FirstName: "Max"},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Printf("Send request: %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
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
				log.Printf("Error while receiving %v\n", err)
				break
			}

			log.Printf("Received: %v\n", res.Result)
		}

		close(waitc)
	}()

	<-waitc
}
