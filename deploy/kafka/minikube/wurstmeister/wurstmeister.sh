#!/bin/bash

# https://stackoverflow.com/questions/44651219/kafka-deployment-on-minikube

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Start Minikube with sufficient resources
minikube start --cpus=2 --memory=4096mb

# Create kubernetes namespace for kafka
kubectl create namespace kafka --dry-run=client -o yaml | kubectl apply -f -

# Deploy Zookeeper
kubectl apply -f zookeeper-deployment.yaml --namespace kafka
kubectl apply -f zookeeper-service.yaml --namespace kafka

# Deploy Kafka (substitute Minikube IP into advertised listeners for external access)
MINIKUBE_IP=$(minikube ip)
sed "s/\${MINIKUBE_IP}/$MINIKUBE_IP/g" kafka-deployment.yaml | kubectl apply -f - --namespace kafka
kubectl apply -f kafka-service.yaml --namespace kafka

# Verify the Deployment
kubectl get pods --namespace kafka
kubectl get services --namespace kafka

echo ""
echo "Kafka is exposed from outside Minikube at: ${MINIKUBE_IP}:30093"
echo ""
echo "Connect from local host:"
echo "  kcat -P -b ${MINIKUBE_IP}:30093 -t test-topic"
echo "  kcat -C -b ${MINIKUBE_IP}:30093 -t test-topic"
echo ""
echo "Or use: minikube service kafka --namespace kafka --url"