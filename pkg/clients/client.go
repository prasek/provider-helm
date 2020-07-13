package clients

import (
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewRestConfig returns a rest config given a secret with connection information.
func NewRestConfig(s *corev1.Secret) (*rest.Config, error) {
	u, err := url.Parse(string(s.Data[runtimev1alpha1.ResourceCredentialsSecretEndpointKey]))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse Kubernetes endpoint as URL")
	}

	return &rest.Config{
		Host:     u.String(),
		Username: string(s.Data[runtimev1alpha1.ResourceCredentialsSecretUserKey]),
		Password: string(s.Data[runtimev1alpha1.ResourceCredentialsSecretPasswordKey]),
		TLSClientConfig: rest.TLSClientConfig{
			// This field's godoc claims clients will use 'the hostname used to
			// contact the server' when it is left unset. In practice clients
			// appear to use the URL, including scheme and port.
			ServerName: u.Hostname(),
			CAData:     s.Data[runtimev1alpha1.ResourceCredentialsSecretCAKey],
			CertData:   s.Data[runtimev1alpha1.ResourceCredentialsSecretClientCertKey],
			KeyData:    s.Data[runtimev1alpha1.ResourceCredentialsSecretClientKeyKey],
		},
		BearerToken: string(s.Data[runtimev1alpha1.ResourceCredentialsSecretTokenKey]),
	}, nil
}

// NewKubeClient returns a kubernetes client given a secret with connection
// information.
func NewKubeClient(config *rest.Config) (client.Client, error) {
	kc, err := client.New(config, client.Options{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create Kubernetes client")
	}

	return kc, nil
}