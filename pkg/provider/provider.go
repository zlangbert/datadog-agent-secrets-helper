package provider

import "github.com/zlangbert/datadog-secrets-provider-aws-secretsmanager/pkg/secret"

type SecretProvider interface {
	Resolve(handles []string) (results map[string]secret.Result, err error)
}
