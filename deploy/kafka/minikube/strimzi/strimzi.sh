#!/bin/bash

# https://saedhasan.medium.com/setting-up-kafka-on-minikube-k8s-using-strimzi-5cac7870d943

# Start Minikube with sufficient resources
minikube start --cpus=2 --memory=4096mb 

# create kubernetes namespace for kafka 
kubectl create namespace kafka

# Install Kafka using Helm
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' --namespace kafka

# Verify the Deployment
kubectl get pods --namespace kafka
kubectl get services --namespace kafka
kubectl logs deployment/strimzi-cluster-operator --namespace kafka -f

# Deploy Kafka Cluster using Strimzi
# kafka-cluster.yaml is modified for service per each broker 
#.  https://stackoverflow.com/questions/77480906/how-to-access-strimzi-kafka-cluster-running-on-minikube-publically
kubectl apply -f kafka-cluster.yaml --n kafka

# Verify the Kafka Cluster Deployment 
kubectl get kafka --namespace kafka

# port forward to access Kafka Brokers in minikube 
kubectl port-forward service/my-cluster-kafka-0 19094:19094 --namespace kafka
kubectl port-forward service/my-cluster-kafka-1 19095:19095 --namespace kafka
kubectl port-forward service/my-cluster-kafka-2 19096:19096 --namespace kafka
