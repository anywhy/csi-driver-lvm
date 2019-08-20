package lvm

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

type driver struct {
	csiDriver *csicommon.CSIDriver

	endpoint         string
	idServer         *identityServer
	nodeServer       csi.NodeServer
	controllerServer *controllerServer
}

const (
	driverName = "lvm.csi.anywhy.github"
	csiVersion = "1.0.0"
)

// NewDriver new driver instance
func NewDriver(nodeID string, endpoint string) *driver {
	driver := &driver{
		csiDriver: csicommon.NewCSIDriver(driverName, csiVersion, nodeID),
		endpoint:  endpoint,
	}

	driver.csiDriver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	})
	driver.csiDriver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
	})

	// Create GRPC servers
	driver.idServer = NewIdentityServer(driver.csiDriver)
	driver.controllerServer = NewControllerServer(driver.csiDriver)
	driver.nodeServer = NewNodeServer(driver.csiDriver, nodeID)
	return driver
}

func (d *driver) Run() {
	glog.Infof("CSI Driver: %v ", driverName)

	server := csicommon.NewNonBlockingGRPCServer()
	server.Start(d.endpoint, d.idServer, d.controllerServer, d.nodeServer)
	server.Wait()
}
