package provider

import "github.com/zlangbert/datadog-secrets-provider-aws-secretsmanager/pkg/secret"

type SecretProvider interface {
	Resolve(handles []*secret.Handle) (results map[string]secret.Result, err error)
}
