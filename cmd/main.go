package main

import (
	"flag"
	"os"

	"github.com/anywhy/csi-driver-lvm/pkg/lvm"
)

var (
	endpoint   = flag.String("endpoint", "unix://tmp/lvm.sock", "CSI endpoint")
	nodeID     = flag.String("nodeid", "", "node id")
	kubeconfig = flag.String("kubeconfig", "", "required only when running out of cluster.")
)

func main() {
	flag.Parse()
	lvm.Kubeconfig = *kubeconfig

	// start driver
	driver := lvm.NewDriver(*nodeID, *endpoint)
	driver.Run()

	os.Exit(0)
}
