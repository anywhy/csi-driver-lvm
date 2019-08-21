package lvm

import (
	"context"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/util/mount"
)

const (
	// LvmVgName lvm vg name
	LvmVgName = "vgName"
	// DefaultFsType default fs type
	DefaultFsType = "ext4"
)

type nodeServer struct {
	*csicommon.DefaultNodeServer
	client  kubernetes.Interface
	mounter mount.SafeFormatAndMount

	nodeID string
}

var (
	// MasterURL  master endpoint
	MasterURL string
	// Kubeconfig config address
	Kubeconfig string
)

func (ns *nodeServer) GetNodeID() string {
	return ns.nodeID
}

func NewNodeServer(d *csicommon.CSIDriver, nodeID string) csi.NodeServer {
	cfg, err := clusterConfig()
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
		client:            kubeClient,
		nodeID:            nodeID,
		mounter: mount.SafeFormatAndMount{
			Interface: mount.New(""),
			Exec:      mount.NewOsExec(),
		},
	}
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	isMnt, err := ns.mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			return nil, status.Error(codes.NotFound, "TargetPath not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !isMnt {
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}

	// Unmount point
	if err := ns.mounter.Unmount(targetPath); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	glog.Warningf("Debug: %v", req)
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) createVolume(ctx context.Context, volumeID string, vgName string) error {
	_, err := ns.client.CoreV1().PersistentVolumes().Get(volumeID, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// pvQuantity := pv.Spec.Capacity["storage"]
	// pvSize := pvQuantity.Value()
	// pvSize = pvSize / (1024 * 1024 * 1024)
	return nil
}

// clusterConfig get k8s cluster config
func clusterConfig() (*rest.Config, error) {
	if len(MasterURL) > 0 || len(Kubeconfig) > 0 {
		return clientcmd.BuildConfigFromFlags(MasterURL, Kubeconfig)
	}
	return rest.InClusterConfig()
}
