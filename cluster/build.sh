#!/usr/bin/bash -ex

plugin=$1

source ./cluster/helpers.sh

registry=localhost:$registry_port
docker build -t ${registry}/device-plugin-$plugin:latest ./cmd/$(echo $plugin | sed 's/-/\//')
