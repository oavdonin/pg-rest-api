#!/usr/bin/env bash
set -euo pipefail

# DOCKER_KEY_B64=ZjkzYmQ3N2ItODYyMC00YmI3LWJmNmQtYzU0Njk0YjRlNjRk
# DOCKER_USERNAME_B64=YXZkb25pbg==

# docker build --target=app -t registry.hub.docker.com/avdonin/pgapi:latest .
# USERNAME=$(echo -n $DOCKER_USERNAME_B64 | base64 -d)
# PASSWORD=$(echo -n $DOCKER_KEY_B64 | base64 -d)

# docker login registry.hub.docker.com --username $USERNAME --password $PASSWORD
# docker push registry.hub.docker.com/avdonin/pgapi:latest

LOCAL_REPOSITORY_PORT=$(docker container port minikube | grep '5000/tcp' | cut -d: -f2)
docker build --target=app -t 127.0.0.1:${LOCAL_REPOSITORY_PORT}/pgapi:lastest .
docker push 127.0.0.1:${LOCAL_REPOSITORY_PORT}/pgapi:latest

kubectl apply -f k8s/namespace.yml
kubectl apply -f k8s/secrets.yml
kubectl apply -f k8s/deploy-postgres.yml
kubectl rollout status statefulset/postgres -n pgapi
kubectl apply -f k8s/deploy-api.yml
kubectl wait --for=condition=available --timeout=120s deployment/pgapi -n pgapi
echo "API Server is ready at http://localhost:8080/people"
kubectl port-forward service/pgapi-svc -n pgapi 8080:8080




