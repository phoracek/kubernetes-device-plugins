#!/usr/bin/bash -e

source ./cluster/helpers.sh

$gocli run --random-ports --nodes 2 --background kubevirtci/k8s-1.10.3

# refresh with new port numbers
source ./cluster/helpers.sh

$gocli scp /etc/kubernetes/admin.conf - > ./kubeconfig
kubectl --kubeconfig=./kubeconfig config set-cluster kubernetes --server=https://127.0.0.1:$k8s_port
kubectl --kubeconfig=./kubeconfig config set-cluster kubernetes --insecure-skip-tls-verify=true
