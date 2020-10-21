# Running E2E Locally

The E2E tests performs a regression test for the various vulnerability providers.
Basic test motion includes:
1. Deploying an image to a registry, private or public, depends on the kind of test you'd like to run
2. For private registries the e2e will need image pull credentials
3. The test creates a test namespace, with a pod that runs the image from #1
4. The test scans the created test namespace by leveraging the vulnerability data available from the provider

* Create local KIND cluster
    ```shell script
    make create-kind-cluster
    ```
* Export Vulnerability Providers configuration
```shell script
export E2E_API_CONFIG='' && read -r -d '' E2E_API_CONFIG << EOM
providers:
  - kind: "gcr"
    repository: "gcr.io/yourproject"
    creds:
      gcr: |
        {
          "type": "service_account",
          "project_id": "yourproject",
          "private_key_id": "XXX",
          "private_key": "",
          "client_email": "imagevulreader@yourproject.iam.gserviceaccount.com",
          "client_id": "666",
          "auth_uri": "https://accounts.google.com/o/oauth2/auth",
          "token_uri": "https://oauth2.googleapis.com/token",
          "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
          "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/imagevulreader%40yourproject.iam.gserviceaccount.com"
        }
  - kind: "ecr"
    repository: "yourawsaccount.dkr.ecr.us-west-2.amazonaws.com/iskan"
    creds:
      ecr:
        accessKeyId: AWSKEY
        secretAccessKey: AWSSECRET
        region: us-west-2
  - kind: "acr"
    repository: "alcide.azurecr.io/iskan"
    creds:
      acr:
        tenantId: mytenantid
        subscriptionId: subscrrptionId
        clientId: clientId
        clientSecret: clientsecret
        cloudName: "AZUREPUBLICCLOUD"
  - kind: "trivy"
    # Use "*" for a capture all images
    repository: "*"
    creds:
      trivy:
        debugMode: false
EOM     
```

* Make sure the pull secrets can be created (see [e2e-pipeline-runner.sh](e2e-pipeline-runner.sh))
    * Export AWS relevant credentials `AWS_SECRET_ACCESS_KEY, AWS_ACCESS_KEY_ID, AWS_REGION`
    * Export relevant GCR credentials `E2E_GCR_PULLSECRET`
    * Export relevant ACR credentials `AZURE_ACR_SP_USER, AZURE_ACR_SP_PASS`
    
* If you are using local inline scan - make sure to install Trivy

* Run e2e by running: `make e2e`

# Adding New Test

- Add your test spec (make sure to give it a unique tag in the **Describe**)
    ```go
    package tests
    
    import (
        "fmt"
        "github.com/alcideio/iskan/e2e/framework"
        . "github.com/onsi/ginkgo"
    )
    
    var _ = Describe("[sanity][mynewtest] My New Test Description", func() {
        f, _ := framework.NewDefaultFramework("mynewtest")
        Context("mynewtest", func() {
            ...
        })
    })
    ```
- Build the e2e framework `make e2e-build`
- Run your specific **Spec** `bin/e2e.test -v 8 -ginkgo.v -ginkgo.focus="\[mynewtest\]"`