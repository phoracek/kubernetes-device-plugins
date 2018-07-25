#!/usr/bin/bash -e

source ./cluster/helpers.sh

$gocli run --random-ports --nodes 2 --background kubevirtci/k8s-1.10.3
$gocli scp /etc/kubernetes/admin.conf - > ./kubeconfig
