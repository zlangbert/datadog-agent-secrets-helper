package provider

import (
	"github.com/pkg/errors"
	"github.com/zlangbert/datadog-agent-secrets-helper/pkg/config"
	"github.com/zlangbert/datadog-agent-secrets-helper/pkg/secret"
)

// SecretProvider common provider interface
type SecretProvider interface {
	Resolve(handles []*secret.Handle) (results map[string]secret.Result)
}

// GetProvider returns a instantiated provider based on id
func GetProvider(cfg *config.HelperConfig, id string) (p SecretProvider, err error) {

	switch id {
	case "aws-sm":
		p, err = NewAwsSecretsManagerProvider(cfg)
	case "kube-secret":
		p, err = NewKubeSecretsProvider(cfg)
	default:
		err = errors.Errorf("unknown provider: %s", id)
	}

	return p, err
}
