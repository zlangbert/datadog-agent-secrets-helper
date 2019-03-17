package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var (
	handlePattern = regexp.MustCompile(`^(?P<provider>[\w-]+):(?P<id>.+):(?P<key>.+)$`)
)

// input from the Datadog agent
type secretsPayload struct {
	Version string   `json:"version"`
	Handles []string `json:"secrets"`
}

type secretHandle struct {
	provider string
	id       string
	key      string
}

func main() {

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

	// build secrets manager client
	sess := session.Must(session.NewSession())
	manager := secretsmanager.New(sess)

	//
	res := map[string]map[string]string{}
	for _, h := range secrets.Handles {

		// parse handle
		handle, err := parseHandle(h)
		if err != nil {
			res[h] = map[string]string{
				"error": fmt.Sprintf("error parsing secret handle: %v", err),
			}
			continue
		}

		// get secret from secrets manager
		secret, err := manager.GetSecretValue(&secretsmanager.GetSecretValueInput{
			SecretId: &handle.id,
		})
		if err != nil {
			res[h] = map[string]string{
				"error": fmt.Sprintf("error getting secret: %v", err),
			}
			continue
		}

		// get secret payload
		var secretPayload map[string]string
		err = json.Unmarshal([]byte(*secret.SecretString), &secretPayload)
		if err != nil {
			res[h] = map[string]string{
				"error": fmt.Sprintf("error parsing secret value, expected object with string key/values: %v", err),
			}
			continue
		}

		// get requested key from secret data
		value := secretPayload[handle.key]
		if value == "" {
			res[h] = map[string]string{
				"error": fmt.Sprintf("secret value for key '%s' is missing or blank", handle.key),
			}
			continue
		}

		res[h] = map[string]string{
			"value": value,
		}
	}

	output, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("could not serialize result: %s", err)
	}
	fmt.Printf(string(output))
}

func parseHandle(h string) (handle *secretHandle, err error) {

	// extract parts from raw handle
	match := handlePattern.FindStringSubmatch(h)
	if match == nil {
		return nil, fmt.Errorf("unexpected handle format: %s", h)
	}

	// build secretHandle
	handle = &secretHandle{
		provider: match[1],
		id:       match[2],
		key:      match[3],
	}

	// validate provider
	if handle.provider != "aws-sm" {
		return nil, fmt.Errorf("unexpected provider in handle: %s", h)
	}

	return handle, nil
}
