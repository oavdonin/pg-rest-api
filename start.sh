#!/usr/bin/env bash
set -euo pipefail

DOCKER_KEY_B64=ZjkzYmQ3N2ItODYyMC00YmI3LWJmNmQtYzU0Njk0YjRlNjRk
DOCKER_USERNAME_B64=YXZkb25pbg==

docker build --target=app -t registry.hub.docker.com/avdonin/pgapi:latest .
USERNAME=$(echo -n $DOCKER_USERNAME_B64 | base64 -d)
PASSWORD=$(echo -n $DOCKER_KEY_B64 | base64 -d)

docker login registry.hub.docker.com --username $USERNAME --password $PASSWORD
docker push registry.hub.docker.com/avdonin/pgapi:latest

kubectl apply -f k8s/namespace.yml
kubectl apply -f k8s/secrets.yml
kubectl apply -f k8s/deploy-postgres.yml 

