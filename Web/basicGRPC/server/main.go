package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hairizuanbinnoorazman/basic-grpc/ticketing"

	"google.golang.org/grpc"
)

var podID string = "test"

type actualCustomerControllerServer struct {
	ticketing.UnimplementedCustomerControllerServer
}

func (a actualCustomerControllerServer) GetCustomer(context.Context, *ticketing.GetCustomerRequest) (*ticketing.Customer, error) {
	log.Println("Hit Get Customer rpc call")
	defer log.Println("End Get Customer rpc call")
	return &ticketing.Customer{
		Id:        podID,
		FirstName: "acac",
		LastName:  "accqqq",
	}, nil
}

func main() {
	fmt.Println("Server Start")

	var exists bool
	podID, exists = os.LookupEnv("POD_NAME")
	if !exists {
		fmt.Println("Value of podID is test")
	}

	lis, _ := net.Listen("tcp", fmt.Sprintf("0.0.0.0:12345"))
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	ticketing.RegisterCustomerControllerServer(grpcServer, actualCustomerControllerServer{})
	grpcServer.Serve(lis)
}
