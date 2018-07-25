#!/usr/bin/bash -e

node=$1

source ./cluster/helpers.sh

$gocli ssh $node
