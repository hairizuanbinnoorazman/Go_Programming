package main

import (
	"context"
	"fmt"
	"time"

	compute "google.golang.org/api/compute/v1"
)

var GCP_PROJECT_ID = "XXXX"

func main() {
	// Added comment
	fmt.Println("Test")
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		panic("Unable to create compute service")
	}
	for {
		request := computeService.Instances.List(GCP_PROJECT_ID, "us-central1-a")
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
