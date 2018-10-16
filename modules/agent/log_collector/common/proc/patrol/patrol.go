package patrol

import (
	"os"
	"runtime"
	"time"

	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/dlog"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/proc/metric"
)

func PatrolLoop() {
	go func() {
		for {
			time.Sleep(time.Second * 10)

			nowMemUsedMB := getMemUsedMB()
			metric.MetricMem(int64(nowMemUsedMB))
			rate := (nowMemUsedMB * 100) / uint64(g.Config().MaxMemMB)

			dlog.Infof("agent mem used : %dMB, percent : %d%%", nowMemUsedMB, rate)
			//若超50%限制，打印warning
			//超过100%，就退出了
			if rate > 50 {
				dlog.Warningf("your log-agent heap memory used rate, current: %d%%", rate)
			}
			if rate > 100 {
				// 堆内存已超过限制，退出进程
				dlog.Errorf("heap memory size over limit. quit process.[used:%dMB][limit:%dMB][rate:%d]", nowMemUsedMB, g.Config().MaxMemMB, rate)
				os.Exit(1)
			}
		}
	}()
}

func getMemUsedMB() uint64 {
	var sts runtime.MemStats
	runtime.ReadMemStats(&sts)
	// 这里取了mem.Alloc
	ret := sts.HeapAlloc / 1024 / 1024
	return ret
}
