#!/bin/bash
k3d cluster delete cluster-k3d
k3d cluster create cluster-k3d -p "8082:30080@agent[0]" --agents 1 --servers 1 --api-port 6443
kubectl apply -f .
