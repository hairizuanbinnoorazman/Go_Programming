package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		fmt.Printf("unable to load up aws configuration :: %v", err)
		os.Exit(1)
	}
	s3Client := s3.NewFromConfig(awsCfg)
	raw, err := os.ReadFile("plastic.png")
	if err != nil {
		fmt.Printf("cannot read file :: %v", err)
		os.Exit(1)
	}
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("tw-auditeye-staging"),
		Key:    aws.String("6cdc3eea-930a-4e5e-8ac8-fae2bdf58e79.jpg"),
		Body:   bytes.NewReader(raw),
	})
	if err != nil {
		fmt.Printf("cannot upload file :: %v", err)
		os.Exit(1)
	}
}
