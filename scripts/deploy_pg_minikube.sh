#!/bin/bash

# https://www.google.com/search?q=how+to+host+postgres+in+minikube

# Start Minikube with sufficient resources
minikube start --cpus 4 --memory 5g

# Add the Bitnami Helm Repository 
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Install the PostgreSQL Helm chart 
helm install pg-minikube --set auth.postgresPassword=<your-strong-password> bitnami/postgresql

# Verify the deployment 
kubectl get pods -n default
kubectl get services -n default 