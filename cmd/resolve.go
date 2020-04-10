package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zlangbert/datadog-secrets-provider-aws-secretsmanager/pkg/config"
	"github.com/zlangbert/datadog-secrets-provider-aws-secretsmanager/pkg/provider"
	"github.com/zlangbert/datadog-secrets-provider-aws-secretsmanager/pkg/secret"
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

			config := &config.HelperConfig{}

			resolve(config)
		},
	}

	err := resolve.Execute()
	if err != nil {
		log.Fatalf("error parsing arguments: %s", err)
	}
}

func resolve(config *config.HelperConfig) {

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
	secretProvider, err := provider.NewAwsSecretsManagerProvider(config)
	if err != nil {
		log.Fatalf("error initializing provider: %s", err)
	}

	// parse handles
	handles := []*secret.Handle{}
	parsingErrors := map[string]secret.Result{}
	for _, h := range secrets.Handles {
		handle, err := secret.ParseHandle(h)
		if err != nil {
			parsingErrors[h] = secret.Result{
				Error: fmt.Sprintf("error parsing secret handle: %v", err),
			}
			continue
		}
		handles = append(handles, handle)
	}

	// resolve handles
	results, err := secretProvider.Resolve(handles)
	if err != nil {
		log.Fatalf("error resolving secrets: %s", err)
	}

	// merge results with any parsing errors
	for k, v := range parsingErrors {
		results[k] = v
	}

	// write result to stdout
	output, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("could not serialize result: %s", err)
	}
	fmt.Printf(string(output))
}
