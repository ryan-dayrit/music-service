#!/bin/bash

# https://stackoverflow.com/questions/44651219/kafka-deployment-on-minikube

# Start Minikube with sufficient resources
minikube start --cpus=2 --memory=4096mb 

# Create kubernetes namespace for kafka 
kubectl create namespace kafka

# Deploy Zookeeper 
kubectl apply -f zookeeper-deployment.yaml --namespace kafka
kubectl apply -f zookeeper-service.yaml --namespace kafka

# Deploy Kafka 
kubectl apply -f kafka-deployment.yaml --namespace kafka
kubectl apply -f kafka-service.yaml --namespace kafka

# Verify the Deployment
kubectl get pods --namespace kafka
kubectl get services --namespace kafka

# Port Forward Pod to Access Kafka Brokers in minikube
kubectl get pods --namespace kafka 
kubectl port-forward --namespace kafka pod/<POD_NAME> 9093:9093

# Send Message to Kafka Topic
echo "Am I receiving this message?" | kcat -P -b localhost:9093 -t test-topic

# Receive Message from Kafka Topic
kcat -C -b localhost:9093 -t test-topic