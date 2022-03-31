package argocd

import (
	"context"
	"io"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	v1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
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
	request := applicationpkg.ApplicationUpdateSpecRequest{
		Name: &appName,
		Spec: v1alpha1.ApplicationSpec{
			Info: []v1alpha1.Info{{
				Name:  "image.tag",
				Value: tag,
			}},
		},
	}

	_, err := a.client.UpdateSpec(context.Background(), &request)
	if err != nil {
		return err
	}

	defer argoio.Close(a.connection)

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
