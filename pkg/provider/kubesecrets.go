package provider

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zlangbert/datadog-agent-secrets-helper/pkg/config"
	"github.com/zlangbert/datadog-agent-secrets-helper/pkg/secret"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"regexp"
)

var (
	idPattern = regexp.MustCompile(`^(?P<namespace>.+)/(?P<name>.+)$`)
)

type KubeSecretsProvider struct {
	client *kubernetes.Clientset
}

func NewKubeSecretsProvider(cfg *config.HelperConfig) (provider SecretProvider, err error) {

	// build client, first try kubeconfig, then in-cluster
	var c *rest.Config

	// look for local kubeconfig
	if h := os.Getenv("HOME"); h != "" {
		path := filepath.Join(h, ".kube", "config")
		info, _ := os.Stat(path)
		if info != nil {
			c, err = clientcmd.BuildConfigFromFlags("", path)
		}
	}

	// try in-cluster config
	if c == nil {
		c, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed building kubernetes client config")
	}

	client, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed building kubernetes client")
	}

	provider = &KubeSecretsProvider{
		client: client,
	}

	return provider, nil
}

func (provider *KubeSecretsProvider) Resolve(handles []*secret.Handle) (results map[string]secret.Result) {

	// TODO: If multiple keys are desired under one secret, that secret is retrieved multiple times. Optimize by
	//  getting each secret only once
	// evaluate each handle
	secretResults := map[string]secret.Result{}
	for _, handle := range handles {

		// extract parts
		match := idPattern.FindStringSubmatch(handle.Id)
		if match == nil {
			secretResults[handle.Handle] = secret.Result{
				Error: fmt.Sprintf("unexpected handle id format: %v", handle.Id),
			}
			continue
		}

		var (
			namespace = match[1]
			name = match[2]
		)

		// get secret
		s, err := provider.client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			secretResults[handle.Handle] = secret.Result{
				Error: fmt.Sprintf("error getting secret: %v", err),
			}
			continue
		}

		// get secret value
		value := string(s.Data[handle.Key])
		if len(value) <= 0 {
			secretResults[handle.Handle] = secret.Result{
				Error: fmt.Sprintf("secret value for key '%s' is missing or blank", handle.Key),
			}
			continue
		}

		secretResults[handle.Handle] = secret.Result{
			Value: value,
		}
	}

	return secretResults
}
