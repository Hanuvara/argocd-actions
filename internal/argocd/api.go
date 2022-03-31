package argocd

import (
	"context"
	"io"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	argoio "github.com/argoproj/argo-cd/v2/util/io"
)

//nolint:lll // go generate is ugly.
//go:generate mockgen -destination=mocks/api_mock.go -package=mocks github.com/omegion/argocd-actions/internal/argocd Interface
// Interface is an interface for API.
type Interface interface {
	Sync(appName string) error
	SetImageTag(appName string, tag string) error
}

// API is struct for ArgoCD api.
type API struct {
	client     applicationpkg.ApplicationServiceClient
	connection io.Closer
}

// APIOptions is options for API.
type APIOptions struct {
	Address  string
	Token    string
	Insecure bool
	ImageTag string
}

// NewAPI creates new API.
func NewAPI(options APIOptions) API {
	clientOptions := argocdclient.ClientOptions{
		ServerAddr: options.Address,
		AuthToken:  options.Token,
		GRPCWeb:    true,
		Insecure:   options.Insecure,
	}

	connection, client := argocdclient.NewClientOrDie(&clientOptions).NewApplicationClientOrDie()

	return API{client: client, connection: connection}
}

// Set application parameter image.tag for a give application
func (a API) SetImageTag(appName string, tag string) error {
	applicationSpec, err := a.client.Get(context.Background(), &applicationpkg.ApplicationQuery{
		Name: &appName,
	})

	if err != nil {
		// Get App Name Failed
		return err
	}

	for index := range applicationSpec.Spec.Source.Helm.Parameters {
		if applicationSpec.Spec.Source.Helm.Parameters[index].Name == "image.tag" {
			applicationSpec.Spec.Source.Helm.Parameters[index].Value = tag
		}
	}
	request := applicationpkg.ApplicationUpdateSpecRequest{
		Name: &appName,
		Spec: applicationSpec.Spec,
	}

	_, err = a.client.UpdateSpec(context.Background(), &request)
	if err != nil {
		return err
	}

	return nil
}

// Sync syncs given application.
func (a API) Sync(appName string) error {
	request := applicationpkg.ApplicationSyncRequest{
		Name:  &appName,
		Prune: true,
	}

	_, err := a.client.Sync(context.Background(), &request)
	if err != nil {
		return err
	}

	defer argoio.Close(a.connection)

	return nil
}
