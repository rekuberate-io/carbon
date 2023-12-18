#!/bin/bash

kind create cluster --name carbon-dev

helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add influxdata https://helm.influxdata.com/
helm repo update

helm upgrade --install metrics-server metrics-server/metrics-server --namespace kube-system
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack -n monitoring --create-namespace
helm upgrade --install influxdb2 influxdata/influxdb2 -n monitoring

