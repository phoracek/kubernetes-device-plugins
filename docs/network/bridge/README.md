# Network Bridge Device Plugin

POC. WIP.

## Build

```
cd cmd/network/bridge
go build
```

## Requirements

[CNI plugins](https://github.com/containernetworking/plugins)
installed in `/opt/cni/`.

Running daemon for CNI IPAM DHCP plugin.

```
/opt/cni/dhcp daemon
```

Bridge `virbr0` available on the host. User must configure DHCP
server running on the bridge.

Device plugin running.

```
./bridge -v 3 -logtostderr
```

## Usage

Create a pod requesting connection to the bridge.

```
apiVersion: v1
kind: Pod
metadata:
  name: demo-pod3
spec:
  containers:
    - name: nginx
      image: nginx:1.7.9
      resources:
        limits:
          bridge.network.mpolednik.github.io/virbr0: 1
```
