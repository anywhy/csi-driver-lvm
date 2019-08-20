package lvm

import (
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

type driver struct {
	csiDriver *csicommon.CSIDriver
	endpoint  string
}
