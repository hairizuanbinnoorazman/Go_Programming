package main

import (
	"context"
	"fmt"
	"os"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("Starting application")

	SERVICE_ACCOUNT_FILE := os.Getenv("SERVICE_ACCOUNT_FILE")
	GCP_PROJECT_ID := os.Getenv("GCP_PROJECT_ID")
	fmt.Printf("Print GCP_PROJECT_ID: %v\n", GCP_PROJECT_ID)
	fmt.Printf("Print SERVICE_ACCOUNT_FILE: %v\n", SERVICE_ACCOUNT_FILE)

	if GCP_PROJECT_ID == "" {
		panic("We need this value to continue operations")
	}

	ctx := context.Background()
	var computeService *compute.InstancesClient
	var err error
	if SERVICE_ACCOUNT_FILE != "" {
		fmt.Println("create computeservice with credentials file")
		computeService, err = compute.NewInstancesRESTClient(ctx, option.WithCredentialsFile(SERVICE_ACCOUNT_FILE))
	} else {
		fmt.Println("create computeservice with workload identity")
		computeService, err = compute.NewInstancesRESTClient(ctx)
	}
	if err != nil {
		fmt.Println(err)
		panic("Unable to create compute service")
	}
	for {
		request := computeService.List(context.TODO(), &computepb.ListInstancesRequest{
			Project: GCP_PROJECT_ID,
			Zone:    "asia-southeast1-a",
		})
		for {
			zz, err := request.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
			}
			fmt.Printf("Server name: %s\n", *zz.Name)
			time.Sleep(2 * time.Second)
		}
	}
}
