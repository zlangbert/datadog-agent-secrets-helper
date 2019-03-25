package provider

import "github.nike.com/zlangb/dd-secrets-provider-secretsmanager/pkg/secret"

type SecretProvider interface {
	Resolve(handles []string) (results map[string]secret.Result, err error)
}
