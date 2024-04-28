package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func ProvideAWSConfig() (aws.Config, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	return sdkConfig, err
}
