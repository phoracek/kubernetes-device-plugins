package network

import (
	"container/list"
	"fmt"
	"os/exec"
	"time"

	"github.com/golang/glog"
	"github.com/mpolednik/kubernetes-device-plugins/pkg/dockerutils"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1alpha"
)

const (
	hostDevicePath = "/tmp/deviceplugin-network-bridge-fakedev"
)

type AttachmentRequest struct {
	DeviceID      string
	ContainerPath string
}

type AttachmentCallbackType func(deviceClass, deviceID, containerID, containerNetNS string) error

// TODO: use os call
func CreateAttachmentHostDevice() {
	glog.V(3).Info("Creating attachment host device")
	cmd := exec.Command("mknod", hostDevicePath, "b", "1", "1")
	cmd.Run()
}

// TODO: pluginName and deviceClass from ndp
func (ndp *NetworkDevicePlugin) AllocateAttachment(pluginName, deviceID, deviceClass string) *pluginapi.DeviceSpec {
	dev := new(pluginapi.DeviceSpec)
	attachmentPath := getAttachmentPath(pluginName, deviceClass, deviceID)
	dev.HostPath = hostDevicePath
	dev.ContainerPath = attachmentPath
	dev.Permissions = "r"
	ndp.AttachmentCh <- &AttachmentRequest{
		deviceID,
		attachmentPath,
	}
	return dev
}

func getAttachmentPath(pluginName, bridge, nic string) string {
	return fmt.Sprintf("/tmp/%s/%s/%s", pluginName, bridge, nic)
}

func (ndp *NetworkDevicePlugin) RunAttachmentManager() {
	glog.V(3).Info("Running attachBridges procedure")

	pendingAttachments := list.New()

	cli, err := dockerutils.NewClient()
	if err != nil {
		glog.V(3).Info("Failed to connect to Docker")
		panic(err)
	}

	// TODO: Make this more sane and efficient, run CNI in a separate goroutine
	for {
		select {
		case attachmentRequest := <-ndp.AttachmentCh:
			glog.V(3).Infof("Received a new attachment request: %s", attachmentRequest)
			pendingAttachments.PushBack(attachmentRequest)
		case <-ndp.StopCh:
			glog.V(3).Info("Received stop signal")
			return
		default:
			time.Sleep(time.Second)
		}

		for a := pendingAttachments.Front(); a != nil; a = a.Next() {
			attachmentRequest := a.Value.(*AttachmentRequest)
			glog.V(3).Infof("Handling pending attachment request for: %s", attachmentRequest.DeviceID)
			containerID, err := cli.GetContainerIDByMountedDevice(attachmentRequest.ContainerPath)
			if err != nil {
				glog.V(3).Info("Container was not found")
				continue
			}

			containerNetNS, err := cli.GetNetNSByContainerID(containerID)
			if err != nil {
				glog.V(3).Info("Failed to obtain container's netns")
				continue
			}

			err = ndp.AttachmentCallback(ndp.DeviceClass, attachmentRequest.DeviceID, containerID, containerNetNS)
			if err == nil {
				glog.V(3).Info("Successfully attached pod to a bridge")
			} else {
				glog.V(3).Infof("Pod attachment failed with: %s", err)
			}

			pendingAttachments.Remove(a)
		}
	}
}
