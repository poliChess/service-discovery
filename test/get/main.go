package main

import (
	"context"
	"flag"
	"log"
	"service-discovery/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	name = flag.String("n", "", "service name")
)

func main() {
	flag.Parse()

	if *name == "" {
		log.Fatal("specify service name with -n")
	}

	conn, err := grpc.Dial("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("counld not connect to grpc server: %v", err)
	}
	defer conn.Close()

	c := proto.NewServiceDiscoveryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &proto.GetRequest{ServiceName: *name})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}

	log.Println("response:", r)
}
