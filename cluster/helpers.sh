gocli="docker run --net=host --privileged --rm -v /var/run/docker.sock:/var/run/docker.sock kubevirtci/gocli:latest"
gocli_interative="docker run --net=host --privileged --rm -it -v /var/run/docker.sock:/var/run/docker.sock kubevirtci/gocli:latest"
k8s_port=$($gocli ports k8s | tr -d '\r')
registry_port=$($gocli ports registry | tr -d '\r')
