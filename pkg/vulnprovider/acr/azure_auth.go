package acr

import (
	"fmt"
	"github.com/alcideio/iskan/pkg/types"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func getToken(azConfig *types.Azure, cloudEnv *azure.Environment) (*adal.ServicePrincipalToken, error) {
	if azConfig == nil {
		return nil, fmt.Errorf("failed to obtain Azure Configuration")
	}

	oauthConfig, err := adal.NewOAuthConfig(cloudEnv.ActiveDirectoryEndpoint, azConfig.TenantId)
	if err != nil {
		return nil, err
	}

	token, err := adal.NewServicePrincipalToken(*oauthConfig, azConfig.ClientId, azConfig.ClientSecret, cloudEnv.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getCloudEnvironment(azConfig *types.Azure) (azure.Environment, error) {
	if azConfig == nil {
		return azure.Environment{}, fmt.Errorf("failed to obtain Azure Configuration")
	}

	return azure.EnvironmentFromName(azConfig.CloudName)
}

func NewAuthorizer(azConfig *types.Azure) (autorest.Authorizer, error) {
	if azConfig != nil && azConfig.ClientSecret != "" {
		env, err := getCloudEnvironment(azConfig)

		if err != nil {
			return nil, fmt.Errorf("Failed to obtain Azure Cloud Environment")
		}

		token, err := getToken(azConfig, &env)
		if err != nil {
			return nil, fmt.Errorf("Failed to obtain Azure Client Token")
		}

		return autorest.NewBearerAuthorizer(token), nil
	}

	return auth.NewAuthorizerFromCLI()
}
