package log_collector

import (
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/http"

	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/proc/metric"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/proc/patrol"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/utils"

	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/dlog"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/g"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/worker"

	"runtime"
)

func init() {
	g.InitAll()
	defer g.CloseLog()

	maxCoreNum := utils.GetCPULimitNum(g.Conf().MaxCPURate)
	dlog.Infof("bind [%d] cpu core", maxCoreNum)
	runtime.GOMAXPROCS(maxCoreNum)

	go metric.MetricLoop(60)
	go worker.UpdateStrategiesLoop()
	go patrol.PatrolLoop()
	go worker.PusherStart()

	http.Start()
}
