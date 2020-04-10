package provider

import (
	"github.com/pkg/errors"
	"github.com/zlangbert/datadog-secrets-provider-aws-secretsmanager/pkg/config"
	"github.com/zlangbert/datadog-secrets-provider-aws-secretsmanager/pkg/secret"
)

type SecretProvider interface {
	Resolve(handles []*secret.Handle) (results map[string]secret.Result)
}

func GetProvider(cfg *config.HelperConfig, id string) (p SecretProvider, err error)  {

	switch id {
	case "aws-sm":
		p, err = NewAwsSecretsManagerProvider(cfg)
	default:
		err = errors.Errorf("unknown provider: %s", id)
	}

	return p, err
}