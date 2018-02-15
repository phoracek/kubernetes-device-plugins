package bridge

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/mpolednik/kubernetes-device-plugins/pkg/dpm"
	"github.com/mpolednik/kubernetes-device-plugins/pkg/network"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1alpha"
)

const (
	pluginName        = "device-plugin-network-bridge"
	resourceNamespace = "bridge.network.mpolednik.github.io/"
)

type NetworkBridgeDevicePlugin struct {
	network.NetworkDevicePlugin
}

func NewNetworkBridgeDevicePlugin(deviceClass string, devices []*pluginapi.Device) NetworkBridgeDevicePlugin {
	return NetworkBridgeDevicePlugin{
		network.NewNetworkDevicePlugin(pluginName, resourceNamespace, deviceClass, devices, attachmentCallback),
	}
}

// TODO: there is a bug, iface has name from bridge, not slave
func attachmentCallback(deviceClass, deviceID, containerID, containerNetNS string) error {
	pluginConfig := []byte(fmt.Sprintf(`
		{
			"cniVersion": "0.3.1",
			"type": "bridge",
			"name": "%s",
			"bridge": "%s",
			"ipam": {
				"type": "dhcp"
			}
		}
	`, deviceClass, deviceClass))
	err := network.RunCNIPlugin(pluginConfig, deviceID, containerID, containerNetNS)
	return err
}

// TODO: this can be moved to network/plugin.go, so attachment logic is hidden
// TODO: create DevicePluginInitialized func type
func newDevicePlugin(bridge string, nics []string) dpm.DevicePluginInterface {
	var devs []*pluginapi.Device
	for _, nic := range nics {
		devs = append(devs, &pluginapi.Device{
			ID:     nic,
			Health: pluginapi.Healthy,
		})
	}

	glog.V(3).Infof("Creating device plugin %s, initial devices %v", bridge, nics)
	nbdp := NewNetworkBridgeDevicePlugin(bridge, devs)
	ret := &nbdp
	ret.DevicePlugin.Deps = ret

	// TODO: This should be triggered by start()
	network.CreateAttachmentHostDevice()
	go ret.RunAttachmentManager()

	return ret
}
