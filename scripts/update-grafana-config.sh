#!/bin/sh

kubectl delete -f ./deploy/grafana.yaml
kubectl delete configmap grafana-datasources
kubectl create configmap grafana-datasources --from-file=./deploy/grafana-datasources.yaml
kubectl apply -f ./deploy/grafana.yaml
