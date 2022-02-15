#!/bin/bash
kubectl delete -f bookstore-components/
kubectl delete -f monitoring-systems/
kubectl apply -f bookstore-components/
kubectl apply -f monitoring-systems/