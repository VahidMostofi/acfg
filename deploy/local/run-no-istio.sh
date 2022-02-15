#!/bin/bash
k3d cluster delete cluster-k3d
k3d cluster create cluster-k3d -p "8082:30080@agent[0]" --agents 1 --servers 1 --api-port 6443
kubectl apply -f bookstore-components/
kubectl apply -f monitoring-systems/metrics-server-components.yml
kubectl apply -f monitoring-systems/influxdb-secrets.yml
kubectl apply -f monitoring-systems/endpoint-gateway.yml
kubectl apply -f monitoring-systems/k8s-monitor.yml
