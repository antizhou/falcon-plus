package collector

import (
	"github.com/open-falcon/falcon-plus/modules/agent/collector/common/proc/metric"
	"github.com/open-falcon/falcon-plus/modules/agent/collector/common/g"
	"github.com/open-falcon/falcon-plus/modules/agent/collector/worker"
)

func Start() {
	g.InitAll()

	go metric.MetricLoop(60)
	go worker.UpdateStrategiesLoop()
	// memory control
	//go patrol.PatrolLoop()
	go worker.PusherStart()
}
