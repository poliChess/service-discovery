package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"service-discovery/proto"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedServiceDiscoveryServer

	mutex sync.Mutex
	services map[string][]string;
}

func (s *server) Register(ctx context.Context, service *proto.Service) (*proto.Status, error) {
	log.Println("Registering service", service.GetServiceName(), "at", service.GetServiceAddr())

	s.mutex.Lock()
	l := s.services[service.GetServiceName()]
	l = append(l, service.GetServiceAddr())
	s.services[service.GetServiceName()] = l
	s.mutex.Unlock()

	return &proto.Status{
		Success: true,
		Message: "ok",
	}, nil
}

func (s *server) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	name := req.GetServiceName()

	s.mutex.Lock()
	l := s.services[name]

	if len(l) == 0 {
		s.mutex.Unlock()
		return &proto.GetResponse{
			Status: &proto.Status{
				Success: false,
				Message: "service not found",
			},
		}, nil
	}

	addr := l[0]
	l = append(l[1:], addr)
	s.services[name] = l
	s.mutex.Unlock()

	return &proto.GetResponse{
		Service: &proto.Service{
			ServiceName: name,
			ServiceAddr: addr,
		},
		Status: &proto.Status{
			Success: true,
			Message: "ok",
		},
	}, nil
}

var port = 3000

func main() {
	log.Println("starting discovery service")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterServiceDiscoveryServer(s, &server{
		mutex: sync.Mutex{},
		services: make(map[string][]string),
	})

	log.Println("discovery service started")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
