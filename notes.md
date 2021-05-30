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
iskan cluster --cluster-context mycluster --api-config myconfig.yaml

# Get vulnerability information for a specific image
iskan image --image="gcr.io/myproj/path/to/myimage:v1.0" --api-config myconfig.yaml -f table --filter-severity CRITICAL,HIGH
```

##### Vulnerabilities Provider API References

| Provider | References            |
|----------|-----------------------|
| **ECR** | [ECR Policies](https://docs.aws.amazon.com/AmazonECR/latest/userguide/ecr_managed_policies.html), [ECR Image Scanning](https://docs.aws.amazon.com/AmazonECR/latest/userguide/image-scanning.html#describe-scan-findings) |
| **GCR** | [Enabling the Container Scanning API](https://cloud.google.com/container-registry/docs/enabling-disabling-container-analysis#enable-scanning)                      |
| **ACR** | [Azure Defender](https://docs.microsoft.com/en-us/azure/security-center/defender-for-container-registries-introduction), [Vulnerability Assessment in Azure](https://techcommunity.microsoft.com/t5/azure-security-center/exporting-azure-container-registry-vulnerability-assessment-in/ba-p/1255244)|
| **InsightVM** | [InsightVM Container Security](https://www.rapid7.com/products/insightvm/features/container-security/)|
| **Harbor** | [Harbor Administration](https://goharbor.io/docs/2.1.0/administration/vulnerability-scanning/)|
| **Trivy** | [Trivy on GitHub](https://github.com/aquasecurity/trivy)|


