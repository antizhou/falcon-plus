package g

import (
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/dlog"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func InitLog() error {
	backend, err := dlog.NewFileBackend(g.Config().Log.LogPath)
	if err != nil {
		return err
	} else {
		dlog.SetLogging(g.Config().Log.LogLevel, backend)
		// 日志rotate设置
		backend.Rotate(g.Config().Log.LogRotateNum, uint64(1024*1024*g.Config().Log.LogRotateSize))
		return nil
	}
}

func CloseLog() {
	dlog.Close()
}
