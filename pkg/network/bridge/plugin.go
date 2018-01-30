package bridge

import (
	"container/list"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/containernetworking/cni/libcni"
	"github.com/golang/glog"
	"github.com/mpolednik/kubernetes-device-plugins/pkg/dockerutils"
	"github.com/mpolednik/kubernetes-device-plugins/pkg/dpm"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1alpha"
)

const (
	cniBinaries       = "/opt/cni/bin"
	fakeDevicePath    = "/tmp/deviceplugin-network-bridge-fakedev"
	resourceNamespace = "bridge.network.mpolednik.github.io/"
)

type NetworkBridgeDevicePlugin struct {
	dpm.DevicePlugin
	bridge       string
	attachmentCh chan *Attachment
}

func newDevicePlugin(bridge string, nics []string) *NetworkBridgeDevicePlugin {
	var devs []*pluginapi.Device

	for _, nic := range nics {
		devs = append(devs, &pluginapi.Device{
			ID:     nic,
			Health: pluginapi.Healthy,
		})
	}

	glog.V(3).Infof("Creating device plugin %s, initial devices %v", bridge, nics)
	ret := &NetworkBridgeDevicePlugin{
		dpm.DevicePlugin{
			Socket:       pluginapi.DevicePluginPath + bridge,
			Devs:         devs,
			ResourceName: resourceNamespace + bridge,
			StopCh:       make(chan interface{}),
		},
		bridge,
		make(chan *Attachment),
	}
	ret.DevicePlugin.Deps = ret

	// TODO: This should be triggered by start()
	createFakeDevice()
	go ret.attachBridges()

	return ret
}

func createFakeDevice() {
	glog.V(3).Info("Creating fake block device")
	cmd := exec.Command("mknod", fakeDevicePath, "b", "1", "1")
	cmd.Run()
}

func (nbdp *NetworkBridgeDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	var devs []*pluginapi.Device

	for _, d := range nbdp.DevicePlugin.Devs {
		devs = append(devs, &pluginapi.Device{
			ID:     d.ID,
			Health: pluginapi.Healthy,
		})
	}

	s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})

	for {
		select {
		case <-nbdp.StopCh:
			return nil
		}
	}
}

type Attachment struct {
	DeviceID      string
	ContainerPath string
}

func (nbdp *NetworkBridgeDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	var response pluginapi.AllocateResponse

	for _, nic := range r.DevicesIDs {
		dev := new(pluginapi.DeviceSpec)
		attachmentPath := getAttachmentPath(nbdp.bridge, nic)
		dev.HostPath = fakeDevicePath
		dev.ContainerPath = attachmentPath
		dev.Permissions = "r"
		response.Devices = append(response.Devices, dev)
		nbdp.attachmentCh <- &Attachment{
			nic,
			attachmentPath,
		}
	}

	return &response, nil
}

func getAttachmentPath(bridge string, nic string) string {
	return fmt.Sprintf("/tmp/device-plugin-network-bridge/%s/%s", bridge, nic)
}

func (nbdp *NetworkBridgeDevicePlugin) attachBridges() {
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
		case attachment := <-nbdp.attachmentCh:
			glog.V(3).Infof("Received a new attachment request: %s", attachment)
			pendingAttachments.PushBack(attachment)
		case <-nbdp.StopCh:
			glog.V(3).Info("Received stop signal")
			return
		default:
			time.Sleep(time.Second)
		}

		for a := pendingAttachments.Front(); a != nil; a = a.Next() {
			attachment := a.Value.(*Attachment)
			glog.V(3).Infof("Handling pending attachment request for: %s", attachment.DeviceID)
			containerID, err := cli.GetContainerIDByMountedDevice(attachment.ContainerPath)
			if err != nil {
				glog.V(3).Info("Container was not found")
				continue
			}

			containerNetNS, err := cli.GetNetNSByContainerID(containerID)
			if err != nil {
				glog.V(3).Info("Failed to obtain container's netns")
				continue
			}

			err = attachPodToBridge(nbdp.bridge, attachment.DeviceID, containerID, containerNetNS)
			if err == nil {
				glog.V(3).Info("Successfully attached pod to a bridge")
			} else {
				glog.V(3).Infof("Pod attachment failed with: %s", err)
			}
			pendingAttachments.Remove(a)
		}
	}
}

func attachPodToBridge(bridge string, nic string, containerID string, containerNetNS string) error {
	pluginConfig := []byte(fmt.Sprintf(`
    {
      "cniVersion": "0.3.1",
      "name": "%s",
      "type": "bridge",
      "bridge": "%s",
      "ipam": {
        "type": "dhcp"
      }
    }
	`, bridge, bridge))

	netConfig, err := libcni.ConfFromBytes(pluginConfig)
	if err != nil {
		panic(err)
	}

	cniConfig := libcni.CNIConfig{Path: []string{filepath.Dir(cniBinaries)}}

	runtimeConfig := &libcni.RuntimeConf{
		ContainerID: containerID,
		NetNS:       containerNetNS,
		IfName:      nic,
	}

	_, err = cniConfig.AddNetwork(netConfig, runtimeConfig)

	return err
}
