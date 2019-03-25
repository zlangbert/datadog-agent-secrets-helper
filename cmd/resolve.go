package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/cobra"
	"github.nike.com/zlangb/dd-secrets-provider-secretsmanager/pkg/provider"
	"io/ioutil"
	"log"
	"os"
)

// input from the Datadog agent
type secretsPayload struct {
	Version string   `json:"version"`
	Handles []string `json:"secrets"`
}

func Resolve() {

	var resolve = &cobra.Command{
		Use:   "datadog-secrets-provider-aws-secretsmanager",
		Short: "Datadog agent secrets provider backed by AWS Secrets Manager",
		Run: func(cmd *cobra.Command, args []string) {

			// log.Printf("region: %v", cmd.Flag("region").Value.String())
			region := cmd.Flag("region").Value.String()
			creds := credentials.NewStaticCredentials(
				cmd.Flag("access-key-id").Value.String(),
				cmd.Flag("secret-access-key").Value.String(),
				"",
			)

			resolve(region, creds)
		},
	}

	resolve.PersistentFlags().String("region", "", "aws region e.g 'us-west-2'")
	_ = resolve.MarkPersistentFlagRequired("region")
	resolve.PersistentFlags().String("access-key-id", "", "aws access key id")
	_ = resolve.MarkPersistentFlagRequired("access-key-id")
	resolve.PersistentFlags().String("secret-access-key", "", "aws secret access key")
	_ = resolve.MarkPersistentFlagRequired("secret-access-key")

	err := resolve.Execute()
	if err != nil {
		log.Fatalf("error parsing arguments: %s", err)
	}
}

func resolve(region string, creds *credentials.Credentials) {

	// ensure there is data being sent, otherwise the following read could hang waiting for input
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		log.Fatal("secrets payload must be passed through stdin")
	}

	// read and unmarshal input from agent
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("could not read from stdin: %s", err)
	}

	secrets := secretsPayload{}
	err = json.Unmarshal(data, &secrets)
	if err != nil {
		log.Fatalf("could not deseralize input: %s", err)
	}

	// build provider
	secretProvider, err := provider.NewAwsSecretsManagerProvider(
		region,
		creds,
	)
	if err != nil {
		log.Fatalf("error initializing provider: %s", err)
	}

	// resolve handles
	results, err := secretProvider.Resolve(secrets.Handles)
	if err != nil {
		log.Fatalf("error resolving secrets: %s", err)
	}

	// write result to stdout
	output, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("could not serialize result: %s", err)
	}
	fmt.Printf(string(output))
}
