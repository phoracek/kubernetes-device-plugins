package bridge

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/mpolednik/kubernetes-device-plugins/pkg/dpm"
)

const (
	bridgeName   = "virbr0"
	nicsPoolSize = 100
)

type BridgeLister struct{}

func (pci BridgeLister) Discover() *dpm.DeviceMap {
	var devices = make(dpm.DeviceMap)
	for i := 0; i < nicsPoolSize; i++ {
		devices[bridgeName] = append(devices[bridgeName], fmt.Sprintf("%s-nic%d", bridgeName, i))
	}
	glog.V(3).Infof("Discovered devices: %s", devices)
	return &devices
}

func (pci BridgeLister) NewDevicePlugin(bridge string, nics []string) dpm.DevicePluginInterface {
	return dpm.DevicePluginInterface(newDevicePlugin(bridge, nics))
}
