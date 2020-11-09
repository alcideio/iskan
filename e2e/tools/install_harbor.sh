#!/bin/bash -x

echo "Installing Ingress (Contour)"
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

kubectl patch daemonsets -n projectcontour envoy -p '{"spec":{"template":{"spec":{"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'


echo "Installing Harbor"
helm repo add harbor https://helm.goharbor.io
helm install local-harbor harbor/harbor --wait


# insert/update hosts entry
ip_address="127.0.0.1"
host_name="core.harbor.domain"
# find existing instances in the host file and save the line numbers
matches_in_hosts="$(grep -n $host_name /etc/hosts | cut -f1 -d:)"
host_entry="${ip_address} ${host_name}"

echo "Please enter your password if requested."

if [ ! -z "$matches_in_hosts" ]
then
    echo "Updating existing hosts entry."
    # iterate over the line numbers on which matches were found
    while read -r line_number; do
        # replace the text of each line with the desired host entry
        sudo sed -i '' "${line_number}s/.*/${host_entry} /" /etc/hosts
    done <<< "$matches_in_hosts"
else
    echo "Adding new hosts entry."
    echo "$host_entry" | sudo tee -a /etc/hosts > /dev/null
fi

echo "create testing project"
STATUS=$(curl -w '%{http_code}' -H 'Content-Type: application/json' -H 'Accept: application/json' -X POST -u "admin:Harbor12345" -s --insecure "https://$host_name/api/v2.0/projects" --data '{"project_name":"iskan","metadata":{"public":"false"},"storage_limit":-1}')
if [ $STATUS -ne 201 ]; then
		exit 1
fi

# Push Docker images
docker login -u admin -p Harbor12345 $host_name
docker push core.harbor.domain/iskan/vuln_alpine:latest

curl -v -u "admin:Harbor12345" -X POST "https://$host_name/api/v2.0/projects/iskan/repositories/vuln_alpine/artifacts/latest/scan" \
  -H "accept: application/json" \
  -H "X-Request-Id: MyCorrelationId"
