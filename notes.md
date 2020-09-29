# iSKan

<img src="https://github.com/alcideio/iskan/raw/master/iskan.png" alt="iskan" width="128"/>

Harness your existing Container Image Vulnerability Scanning information to your Kubernetes Cluster.

## Install

```shell script
curl https://raw.githubusercontent.com/alcideio/iskan/master/download.sh | bash
```

## Command Line Examples

```shell script
# Scan the cluster pointed by the kubeconfig context 'myctx'
iskan cluster --cluster-context mycluster

# Get vulnerability information for a specific image
iskan image --image="gcr.io/myproj/path/to/myimage:v1.0" --api-config myconfig.yaml -f table --filter-severity CRITICAL,HIGH



```