package azuredevops

import (
	"context"
	"fmt"
	"log"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/operations"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
)

// Aggregates all of the underlying clients into a single data
// type. Each client is ready to use and fully configured with the correct
// AzDO PAT/organization
//
// AggregatedClient uses interfaces derived from the underlying client structs to
// allow for mocking to support unit testing of the funcs that invoke the
// Azure DevOps client.
type aggregatedClient struct {
	CoreClient            core.Client
	BuildClient           build.Client
	GitReposClient        git.Client
	GraphClient           graph.Client
	OperationsClient      operations.Client
	ServiceEndpointClient serviceendpoint.Client
	ctx                   context.Context
}

func getAzdoClient(azdoPAT string, organizationURL string) (*aggregatedClient, error) {
	ctx := context.Background()

	if azdoPAT == "" {
		return nil, fmt.Errorf("the personal access token is required")
	}

	if organizationURL == "" {
		return nil, fmt.Errorf("the url of the Azure DevOps is required")
	}

	connection := azuredevops.NewPatConnection(organizationURL, azdoPAT)

	// client for these APIs (includes CRUD for AzDO projects...):
	//	https://docs.microsoft.com/en-us/rest/api/azure/devops/core/?view=azure-devops-rest-5.1
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): core.NewClient failed.")
		return nil, err
	}

	// client for these APIs (includes CRUD for AzDO build pipelines...):
	//	https://docs.microsoft.com/en-us/rest/api/azure/devops/build/?view=azure-devops-rest-5.1
	buildClient, err := build.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): build.NewClient failed.")
		return nil, err
	}

	// client for these APIs (monitor async operations...):
	//	https://docs.microsoft.com/en-us/rest/api/azure/devops/operations/operations?view=azure-devops-rest-5.1
	operationsClient := operations.NewClient(ctx, connection)

	// client for these APIs (includes CRUD for AzDO service endpoints a.k.a. service connections...):
	//  https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1
	serviceEndpointClient, err := serviceendpoint.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): serviceendpoint.NewClient failed.")
		return nil, err
	}

	// client for these APIs:
	//	https://docs.microsoft.com/en-us/rest/api/azure/devops/git/?view=azure-devops-rest-5.1
	gitReposClient, err := git.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): git.NewClient failed.")
		return nil, err
	}

	//  https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/?view=azure-devops-rest-5.1
	graphClient, err := graph.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): graph.NewClient failed.")
		return nil, err
	}

	aggregatedClient := &aggregatedClient{
		CoreClient:            coreClient,
		BuildClient:           buildClient,
		GitReposClient:        gitReposClient,
		GraphClient:           graphClient,
		OperationsClient:      operationsClient,
		ServiceEndpointClient: serviceEndpointClient,
		ctx:                   ctx,
	}

	log.Printf("getAzdoClient(): Created core, build, operations, and serviceendpoint clients successfully!")
	return aggregatedClient, nil
}
