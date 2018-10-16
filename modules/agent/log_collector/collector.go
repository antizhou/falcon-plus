package log_collector

import (
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/proc/metric"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/g"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/worker"
)

func Start() {
	g.InitAll()
	defer g.CloseLog()

	go metric.MetricLoop(60)
	go worker.UpdateStrategiesLoop()
	// memory control
	//go patrol.PatrolLoop()
	go worker.PusherStart()
}
