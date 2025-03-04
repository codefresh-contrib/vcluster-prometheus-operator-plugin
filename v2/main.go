package main

import (
	"github.com/codefresh-contrib/vcluster-prometheus-operator-plugin/v2/syncers/podmonitor"
	"github.com/codefresh-contrib/vcluster-prometheus-operator-plugin/v2/syncers/servicemonitor"
	"github.com/loft-sh/vcluster-sdk/plugin"
)

func main() {
	ctx := plugin.MustInit()
	plugin.MustRegister(servicemonitor.NewServiceMonitorSyncer(ctx))
	plugin.MustRegister(podmonitor.NewPodMonitorSyncer(ctx))
	plugin.MustStart()
}
