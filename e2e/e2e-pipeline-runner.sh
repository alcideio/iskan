#!/bin/bash

retval=''

function install_trivy_if_needed() {
  mkdir -p ~/.iskan/trivy/.cache/reports || true

  if ! command -v trivy &> /dev/null
  then
    echo "Downloading Trivy"
    curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/master/contrib/install.sh | sudo sh -s -- -b /usr/local/bin
    echo "Downloading Vuln DB"
    trivy image --download-db-only
  fi
}

function backup_docker_config() {
  echo "Backup docker config"
  ls -la ~/.docker/
  mv ~/.docker/config.json ~/.docker/config_back.json || true
}

function restore_docker_config() {
  echo "Restore docker config"
  mv ~/.docker/config_back.json ~/.docker/config.json || true
}

function ecr_docker_login() {
    backup_docker_config
    if [ "$E2E_PIPELINE" = "YES" ] #if we are in the pipeline
    then
    #export AWS_SECRET_ACCESS_KEY=MY SECRET
    #export AWS_ACCESS_KEY_ID=MY AWS KEY ID
    #export AWS_REGION=us-west-2
      aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 893825821121.dkr.ecr.us-west-2.amazonaws.com
    else
      aws --profile iskan ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 893825821121.dkr.ecr.us-west-2.amazonaws.com
    fi
    retval=`jq -c . ~/.docker/config.json`
    restore_docker_config

    export E2E_ECR_PULLSECRET=$retval
}

function acr_docker_login() {
    #export AZURE_ACR_SP_USER=
    #export AZURE_ACR_SP_PASS=

    backup_docker_config
    docker login -u $AZURE_ACR_SP_USER -p $AZURE_ACR_SP_PASS alcide.azurecr.io
    retval=`jq -c . ~/.docker/config.json`
    restore_docker_config

    export E2E_ACR_PULLSECRET=$retval
}

install_trivy_if_needed

ecr_docker_login
acr_docker_login

# bin/e2e.test  -v 7 -ginkgo.v -ginkgo.focus="\[acr\]"
bin/e2e.test  -v 7 -ginkgo.v
