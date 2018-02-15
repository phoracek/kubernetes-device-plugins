package network

import (
	"path/filepath"

	"github.com/containernetworking/cni/libcni"
)

const (
	cniBinaries = "/opt/cni/bin"
)

func RunCNIPlugin(pluginConfig []byte, ifName, containerID, containerNetNS string) error {
	netConfig, err := libcni.ConfFromBytes(pluginConfig)
	if err != nil {
		return err
	}

	cniConfig := libcni.CNIConfig{Path: []string{filepath.Dir(cniBinaries)}}

	runtimeConfig := &libcni.RuntimeConf{
		ContainerID: containerID,
		NetNS:       containerNetNS,
		IfName:      ifName,
	}

	_, err = cniConfig.AddNetwork(netConfig, runtimeConfig)

	return err
}
