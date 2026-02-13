#!/bin/bash

# https://www.google.com/search?q=deploying+kafka+to+minikube+using+helm

# Start Minikube with sufficient resources
minikube start --cpus=2 --memory=4096mb 

# Add the Bitnami Helm Repository 
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# create kubernetes namespace for kafka 
kubectl create namespace kafka

# Install Kafka using Helm
helm install kafka bitnami/kafka --namespace kafka

# Verify the Deployment
kubectl get pods --namespace kafka
kubectl get services --namespace kafka

# port forward to access Kafka in minikube 
kubectl port-forward svc/my-kafka-headless 9092:9092 --namespace kafka

# this doesn't work 
# WARNING: Since August 28th, 2025, only a limited subset of images/charts are available for free.
#    Subscribe to Bitnami Secure Images to receive continued support and security updates.
#    More info at https://bitnami.com and https://github.com/bitnami/containers/issues/83267
#
#ryandayrit@Ryans-MBP ~ % kubectl get pods --namespace kafka
#NAME                 READY   STATUS                  RESTARTS   AGE
#kafka-controller-0   0/1     Init:ImagePullBackOff   0          27s
#kafka-controller-1   0/1     Init:ImagePullBackOff   0          27s
#kafka-controller-2   0/1     Init:ImagePullBackOff   0          27s
#
# <- these images aren't free