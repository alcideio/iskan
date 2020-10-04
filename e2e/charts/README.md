# Installing Test Pods for Local Tests

```shell script
helm upgrade -i distroless e2e/charts/image-deployer --set image.repository=gcr.io/dcvisor-162009/iskan/e2e/zerovuln_distroless:latest --set image.pullSecretToken="YOUR BASE64 TOKEN"
```

```shell script
helm upgrade -i scratch e2e/charts/image-deployer --set image.repository=gcr.io/dcvisor-162009/iskan/e2e/zerovuln_scratch:latest --set image.pullSecretToken="YOUR BASE64 TOKEN"
```