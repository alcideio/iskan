#!/bin/bash

function acr_login() {
    #export AZURE_ACR_SP_USER=
    #export AZURE_ACR_SP_PASS=
    if [ -z "$AZURE_ACR_SP_PASS" ]
    then
          echo "\$AZURE_ACR_SP_PASS is empty"
          az acr login --name alcide
    else
          docker login -u $AZURE_ACR_SP_USER -p $AZURE_ACR_SP_PASS alcide.azurecr.io
    fi
}

function ecr_login() {
    #export AWS_SECRET_ACCESS_KEY=MY SECRET
    #export AWS_ACCESS_KEY_ID=MY AWS KEY ID
    #export AWS_REGION=us-west-2
    if [ -z "$AWS_ACCESS_KEY_ID" ]
    then
          echo "\$AWS_ACCESS_KEY_ID is empty"
          aws --profile iskan ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 893825821121.dkr.ecr.us-west-2.amazonaws.com
    else
          aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 893825821121.dkr.ecr.us-west-2.amazonaws.com
    fi
}

function gcr_login() {
    #export GCR_SERVICE_ACCOUNT=
    if [ -z "$GCR_SERVICE_ACCOUNT" ]
    then
          echo "\$GCR_SERVICE_ACCOUNT is empty"
          gcloud auth configure-docker && gcloud auth login
    else
          echo -n $GCR_SERVICE_ACCOUNT | docker login -u _json_key --password-stdin gcr.io/dcvisor-162009
    fi
}

function docker_hub_login() {
    if [ -z "$ALCIDE_DOCKER_HUB_TOKEN" ]
    then
          echo "\$ALCIDE_DOCKER_HUB_TOKEN is empty"
    else
          docker login --username alcide --password $ALCIDE_DOCKER_HUB_TOKEN
    fi
}

docker_hub_login
ecr_login
acr_login
gcr_login