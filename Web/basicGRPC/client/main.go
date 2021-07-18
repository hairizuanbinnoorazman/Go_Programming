package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hairizuanbinnoorazman/basic-grpc/ticketing"

	"google.golang.org/grpc"
)

func main() {
	domain, exists := os.LookupEnv("SERVER_DOMAIN")
	if !exists {
		domain = "localhost"
	}

	port, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		port = "12345"
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(3*time.Second), grpc.WithInsecure())
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", domain, port), opts...)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	for {
		getCustomerDetails(conn)
		time.Sleep(3 * time.Second)
	}
}

func getCustomerDetails(conn *grpc.ClientConn) {
	client := ticketing.NewCustomerControllerClient(conn)
	log.Println("Start GetCustomerDetails")
	defer log.Println("End GetCustomerDetails")
	zz, err := client.GetCustomer(context.Background(), &ticketing.GetCustomerRequest{})
	if err != nil {
		fmt.Println(err)
	}
	log.Println(zz)
}
