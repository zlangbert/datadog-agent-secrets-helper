package provider

import "github.com/aws/aws-sdk-go/aws/credentials"

type AwsConfig struct {
	Region      string
	Credentials *credentials.Credentials
}
