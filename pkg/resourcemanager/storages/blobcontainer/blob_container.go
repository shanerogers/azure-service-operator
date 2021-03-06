// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package blobcontainer

import (
	"context"
	"net/http"

	s "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/azure-service-operator/pkg/resourcemanager/config"
	"github.com/Azure/azure-service-operator/pkg/resourcemanager/iam"
	"github.com/Azure/go-autorest/autorest"
)

type AzureBlobContainerManager struct {
	creds config.Credentials
}

func getContainerClient(creds config.Credentials) (s.BlobContainersClient, error) {
	containersClient := s.NewBlobContainersClientWithBaseURI(config.BaseURI(), creds.SubscriptionID())
	auth, err := iam.GetResourceManagementAuthorizer(creds)
	if err != nil {
		return s.BlobContainersClient{}, err
	}
	containersClient.Authorizer = auth
	containersClient.AddToUserAgent(config.UserAgent())
	return containersClient, nil
}

// Creates a blob container in a storage account.
// Parameters:
// resourceGroupName - name of the resource group within the azure subscription.
// accountName - the name of the storage account
// containerName - the name of the container
// accessLevel - 'PublicAccessContainer', 'PublicAccessBlob', or 'PublicAccessNone'
func (m *AzureBlobContainerManager) CreateBlobContainer(ctx context.Context, resourceGroupName string, accountName string, containerName string, accessLevel s.PublicAccess) (*s.BlobContainer, error) {
	containerClient, err := getContainerClient(m.creds)
	if err != nil {
		return nil, err
	}

	blobContainerProperties := s.ContainerProperties{
		PublicAccess: accessLevel,
	}

	container, err := containerClient.Create(
		ctx,
		resourceGroupName,
		accountName,
		containerName,
		s.BlobContainer{ContainerProperties: &blobContainerProperties})

	if err != nil {
		return nil, err
	}

	return &container, err
}

// Get gets the description of the specified blob container.
// Parameters:
// resourceGroupName - name of the resource group within the azure subscription.
// accountName - the name of the storage account
// containerName - the name of the container
func (m *AzureBlobContainerManager) GetBlobContainer(ctx context.Context, resourceGroupName string, accountName string, containerName string) (result s.BlobContainer, err error) {
	containerClient, err := getContainerClient(m.creds)
	if err != nil {
		return s.BlobContainer{}, err
	}

	return containerClient.Get(ctx, resourceGroupName, accountName, containerName)
}

// Deletes a blob container in a storage account.
// Parameters:
// resourceGroupName - name of the resource group within the azure subscription.
// accountName - the name of the storage account
// containerName - the name of the container
func (m *AzureBlobContainerManager) DeleteBlobContainer(ctx context.Context, resourceGroupName string, accountName string, containerName string) (result autorest.Response, err error) {
	containerClient, err := getContainerClient(m.creds)
	if err != nil {
		return autorest.Response{
			Response: &http.Response{
				StatusCode: 500,
			},
		}, err
	}

	return containerClient.Delete(ctx,
		resourceGroupName,
		accountName,
		containerName)
}
