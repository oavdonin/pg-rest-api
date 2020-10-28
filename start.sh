#!/usr/bin/env bash
set -euo pipefail

eval $(minikube -p minikube docker-env)
docker build --target=app -t pgapi:latest .
kubectl apply -f k8s/namespace.yml
kubectl apply -f k8s/secrets.yml
kubectl apply -f k8s/deploy-postgres.yml
kubectl rollout status statefulset/postgres -n pgapi
kubectl apply -f k8s/deploy-api.yml
kubectl wait --for=condition=available --timeout=120s deployment/pgapi -n pgapi
echo "API Server will be ready at http://localhost:8080/people"
kubectl port-forward service/pgapi-svc -n pgapi 8080:8080




