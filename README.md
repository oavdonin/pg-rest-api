# API-Exercise

Tech Stack:

Golang 1.15
Postgres 10
Minikube v1.14.2 (docker-machine-driver-hyperkit, Kubernetes v1.19)

Environment preparation:

Install Minikube with docker driver (default for Mac, Linux)

Start it: minikube start --memory='4g'

Please ensure that you have minikube activated with the next plugins enabled:

- default-storageclass (default in Mac/brew)
- storage-provisioner (default in Mac/brew)
- registry (you have to "minikube addons enable registry")

You can get the list of activated plugins by running: `minikube addons list`

If you don't have them activated, you can do it by executing:

minikube addons enable ${ADDON_NAME}

Deploying REST-API service:

This setup script was tested with:
- MacOS Catalina 10.15.7
-

In order to build, deploy and start you need to have kubectl context associated with your minikube cluster
docker binary should be in your PATH$
Deploy script supports at least zsh and bash

run ./start.sh to get you REST-API running on http://localhost:8080





