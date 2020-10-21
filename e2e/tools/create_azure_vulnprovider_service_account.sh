#!/bin/bash -x

# https://docs.microsoft.com/en-us/azure/container-registry/container-registry-auth-kubernetes#code-try-0

# Modify for your environment.
# ACR_NAME: The name of your Azure Container Registry
# SERVICE_PRINCIPAL_NAME: Must be unique within your AD tenant
SERVICE_PRINCIPAL_NAME=iskan-e2e-azure-vulprovide-service-principal
SUBSCRIPTION="/subscriptions/9efc9618-47a0-4e98-b31e-7194f25188d4"

az ad sp create-for-rbac --name http://$SERVICE_PRINCIPAL_NAME  --scope $SUBSCRIPTION  --sdk-auth
