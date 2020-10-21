#!/bin/bash -x

# https://docs.microsoft.com/en-us/azure/container-registry/container-registry-auth-kubernetes#code-try-0

# Modify for your environment.
# ACR_NAME: The name of your Azure Container Registry
# SERVICE_PRINCIPAL_NAME: Must be unique within your AD tenant
ACR_NAME=alcide
SERVICE_PRINCIPAL_NAME=iskan-e2e-acr-service-principal
SUBSCRIPTION=9efc9618-47a0-4e98-b31e-7194f25188d4

# Obtain the full registry ID for subsequent command args
ACR_REGISTRY_ID=$(az acr show --subscription $SUBSCRIPTION --name $ACR_NAME --query id --output tsv)

# Create the service principal with rights scoped to the registry.
# Default permissions are for docker pull access. Modify the '--role'
# argument value as desired:
# acrpull:     pull only
# acrpush:     push and pull
# owner:       push, pull, and assign roles
SP_PASSWD=$(az ad sp create-for-rbac --subscription $SUBSCRIPTION --name http://$SERVICE_PRINCIPAL_NAME --scopes $ACR_REGISTRY_ID --role acrpull --query password --output tsv)
SP_APP_ID=$(az ad sp show --subscription $SUBSCRIPTION --id http://$SERVICE_PRINCIPAL_NAME --query appId --output tsv)

# Output the service principal's credentials; use these in your services and
# applications to authenticate to the container registry.
echo "Service principal ID: $SP_APP_ID"
echo "Service principal password: $SP_PASSWD"