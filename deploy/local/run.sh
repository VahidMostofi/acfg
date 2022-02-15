#!/bin/bash
k3d cluster delete cluster-k3d
k3d cluster create cluster-k3d --servers 1 --agents 1 --port 9080:80@loadbalancer --port 9443:443@loadbalancer --port 9099:9099@loadbalancer --api-port 6443 --k3s-server-arg '--no-deploy=traefik'
export PATH=$PATH:$HOME/.istioctl/bin
istioctl install --set profile=default --set values.pilot.env.PILOT_HTTP10=1 -y
kubectl label namespace default istio-injection=enabled
kubectl apply -f bookstore-components/
kubectl apply -f monitoring-systems/metrics-server-components.yml
kubectl apply -f monitoring-systems/influxdb-secrets.yml
kubectl apply -f monitoring-systems/endpoint-gateway.yml
kubectl apply -f monitoring-systems/k8s-monitor.yml