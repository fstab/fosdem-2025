#!/bin/bash

make

kind create cluster
kind load docker-image fosdem-2025/product-service
kind load docker-image fosdem-2025/inventory-service
kind load docker-image fosdem-2025/pricing-service

kubectl create configmap beyla-config --from-file=./deploy/beyla-config.yaml
kubectl create configmap prometheus-config --from-file=./deploy/prometheus-config.yaml
kubectl create configmap load-generator --from-file=./deploy/load-generator.js
kubectl create configmap grafana-provisioning --from-file=./deploy/grafana-datasources.yaml --from-file=./deploy/grafana-dashboards.yaml
kubectl create configmap grafana-dashboards --from-file=./deploy/grafana-dashboard-red-metrics.json
kubectl create configmap tempo-config --from-file=./deploy/tempo-config.yaml

kubectl apply \
	-f ./deploy/inventory-service.yaml \
	-f ./deploy/pricing-service.yaml \
	-f ./deploy/product-service.yaml \
	-f ./deploy/beyla.yaml \
	-f ./deploy/prometheus.yaml \
	-f ./deploy/load-generator.yaml \
	-f ./deploy/grafana.yaml \
	-f ./deploy/tempo.yaml \
