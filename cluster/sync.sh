#!/usr/bin/bash -e

plugin=$1

source ./cluster/helpers.sh

registry=localhost:$registry_port
docker push $registry/device-plugin-$plugin
