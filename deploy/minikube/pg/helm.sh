#!/bin/bash

# https://www.google.com/search?q=how+to+host+postgres+in+minikube

# Start Minikube with sufficient resources
minikube start --cpus 4 --memory 5g

# Add the Bitnami Helm Repository 
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Install the PostgreSQL Helm chart 
helm install pg --set auth.postgresPassword=umagos bitnami/postgresql --namespace pg

# Verify the deployment 
kubectl get pods --namespace pg
kubectl get services --namespace pg

# port forward to access PostgreSQL in minikube 
kubectl port-forward svc/pg-minikube-postgresql-hl 5432:5432 --namespace pg
