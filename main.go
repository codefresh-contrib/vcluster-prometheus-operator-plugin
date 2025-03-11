package main

import (
	"github.com/codefresh-contrib/vcluster-prometheus-operator-plugin/syncers/podmonitor"
	"github.com/codefresh-contrib/vcluster-prometheus-operator-plugin/syncers/servicemonitor"
	"github.com/loft-sh/vcluster-sdk/plugin"
	"k8s.io/klog/v2"
)

func main() {
	ctx := plugin.MustInit()

	serviceMonitorSyncer, err := servicemonitor.NewServiceMonitorSyncer(ctx)
	if err != nil {
		klog.Fatalf("new servicemonitor syncer: %v", err)
	}
	plugin.MustRegister(serviceMonitorSyncer)

	podMonitorSyncer, err := podmonitor.NewPodMonitorSyncer(ctx)
	if err != nil {
		klog.Fatalf("new podmonitor syncer: %v", err)
	}
	plugin.MustRegister(podMonitorSyncer)

	plugin.MustStart()
}
