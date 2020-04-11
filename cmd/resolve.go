package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zlangbert/datadog-agent-secrets-helper/pkg/config"
	"github.com/zlangbert/datadog-agent-secrets-helper/pkg/provider"
	"github.com/zlangbert/datadog-agent-secrets-helper/pkg/secret"
	"io/ioutil"
	"log"
	"os"
)

// input from the Datadog agent
type secretsPayload struct {
	Version string   `json:"version"`
	Handles []string `json:"secrets"`
}

// Resolve parses args options and starts the secret resolution process for this invocation
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

func resolve(cfg *config.HelperConfig) {

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

	// result accumulation
	results := map[string]secret.Result{}

	// parse handles
	handles := []*secret.Handle{}
	for _, h := range secrets.Handles {
		handle, err := secret.ParseHandle(h)
		if err != nil {
			results[h] = secret.Result{
				Error: fmt.Sprintf("error parsing secret handle: %v", err),
			}
			continue
		}

		// append any parsed handles
		handles = append(handles, handle)
	}

	// group handles
	groups := groupHandlesByProvider(handles)

	// resolve handle groups
	for providerID, handles := range groups {
		for h, result := range resolveGroup(cfg, providerID, handles) {
			results[h] = result
		}
	}

	// write result to stdout
	output, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("could not serialize result: %s", err)
	}
	fmt.Printf(string(output))
}

func groupHandlesByProvider(handles []*secret.Handle) map[string][]*secret.Handle {
	groups := map[string][]*secret.Handle{}
	for _, h := range handles {
		if _, ok := groups[h.Provider]; ok {
			// append h to existing list
			groups[h.Provider] = append(groups[h.Provider], h)
		} else {
			// create new group
			groups[h.Provider] = []*secret.Handle{h}
		}
	}

	return groups
}

func resolveGroup(cfg *config.HelperConfig, providerID string, handles []*secret.Handle) map[string]secret.Result {
	results := map[string]secret.Result{}

	p, err := provider.GetProvider(cfg, providerID)
	if err != nil {

		// if we failed to initialize the provider then return that error for all handles
		for _, handle := range handles {
			results[handle.Handle] = secret.Result{
				Error: fmt.Sprintf("error initializing provider: %v", err),
			}
		}
	} else {

		// resolve group with provider
		results = p.Resolve(handles)
	}

	return results
}
