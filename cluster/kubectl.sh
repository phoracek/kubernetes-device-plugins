#!/usr/bin/bash -e

source ./cluster/helpers.sh

kubectl --kubeconfig ./kubeconfig --insecure-skip-tls-verify --server https://localhost:$k8s_port "$@"
