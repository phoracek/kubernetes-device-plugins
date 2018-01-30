package main

import (
	"flag"

	"github.com/mpolednik/kubernetes-device-plugins/pkg/dpm"
	"github.com/mpolednik/kubernetes-device-plugins/pkg/network/bridge"
)

func main() {
	flag.Parse()

	manager := dpm.NewDevicePluginManager(bridge.BridgeLister{})
	manager.Run()
}
