package network

import (
	"context"

	"github.com/mpolednik/kubernetes-device-plugins/pkg/dpm"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1alpha"
)

type NetworkDevicePlugin struct {
	dpm.DevicePlugin
	PluginName         string
	DeviceClass        string
	AttachmentCh       chan *AttachmentRequest
	AttachmentCallback AttachmentCallbackType
}

func NewNetworkDevicePlugin(pluginName, resourceNamespace, deviceClass string, devices []*pluginapi.Device, attachmentCallback AttachmentCallbackType) NetworkDevicePlugin {
	return NetworkDevicePlugin{
		dpm.NewDevicePlugin(resourceNamespace, deviceClass, devices),
		pluginName,
		deviceClass,
		make(chan *AttachmentRequest),
		attachmentCallback,
	}
}

func (ndp *NetworkDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	s.Send(&pluginapi.ListAndWatchResponse{Devices: ndp.DevicePlugin.Devs})

	for {
		select {
		case <-ndp.StopCh:
			return nil
		}
	}
}

func (ndp *NetworkDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	var response pluginapi.AllocateResponse

	for _, deviceID := range r.DevicesIDs {
		dev := ndp.AllocateAttachment(ndp.PluginName, ndp.DeviceClass, deviceID)
		response.Devices = append(response.Devices, dev)
	}

	return &response, nil
}
