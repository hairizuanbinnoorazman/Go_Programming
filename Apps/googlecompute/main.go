package main

import (
	"context"
	"fmt"
	"os"
	"time"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("Starting application")

	SERVICE_ACCOUNT_FILE := os.Getenv("SERVICE_ACCOUNT_FILE")
	GCP_PROJECT_ID := os.Getenv("GCP_PROJECT_ID")
	REGION := os.Getenv("REGION")
	fmt.Printf("Print GCP_PROJECT_ID: %v\n", GCP_PROJECT_ID)

	ctx := context.Background()
	var computeService *compute.Service
	var err error
	if SERVICE_ACCOUNT_FILE != "" {
		fmt.Println("create computeservice with credentials file")
		computeService, err = compute.NewService(ctx, option.WithCredentialsFile(SERVICE_ACCOUNT_FILE))
	} else {
		fmt.Println("create computeservice with workload identity")
		computeService, err = compute.NewService(ctx)
	}
	if err != nil {
		panic("Unable to create compute service")
	}
	for {
		request := computeService.Instances.List(GCP_PROJECT_ID, REGION)
		zz, err := request.Do()
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
		for _, z := range zz.Items {
			fmt.Printf("Server name: %s\n", z.Name)
		}
		time.Sleep(2 * time.Second)
	}

}
