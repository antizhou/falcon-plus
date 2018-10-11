package worker

import (
	"sync"

	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/dlog"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/g"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/scheme"

	"time"

	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/reader"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/strategy"
)

// ConfigInfo to control config
type ConfigInfo struct {
	ID       int64
	FilePath string
}

// Job to control job
type Job struct {
	r *reader.Reader
	w *WorkerGroup
}

// JobManager to manage jobs
var JobManager map[string]*Job //管理job,文件路径为key
// JobManagerLock is a global lock
var JobManagerLock *sync.RWMutex

// ManagerConfig to manage configs
var ManagerConfig map[int64]*ConfigInfo

func init() {
	JobManager = make(map[string]*Job)
	JobManagerLock = new(sync.RWMutex)
	ManagerConfig = make(map[int64]*ConfigInfo)
}

// UpdateStrategiesLoop to update strategys
func UpdateStrategiesLoop() {
	for {
		strategy.Update()
		strategyMap := strategy.GetAll() //最新策略
		JobManagerLock.Lock()

		for id, st := range strategyMap {
			config := &ConfigInfo{
				ID:       id,
				FilePath: st.FilePath,
			}
			cache := make(chan string, g.Conf().Worker.QueueSize)
			if err := createJob(config, cache, st); err != nil {
				dlog.Errorf("create job fail [id:%d][filePath:%s][err:%v]", config.ID, config.FilePath, err)
			}
		}

		for id := range ManagerConfig {
			if _, ok := strategyMap[id]; !ok { //如果策略中不存在，说明用户已删除
				config := &ConfigInfo{
					ID:       id,
					FilePath: ManagerConfig[id].FilePath,
				}
				deleteJob(config)
			}
		}
		JobManagerLock.Unlock()

		//更新counter
		GlobalCount.UpdateByStrategy(strategyMap)
		time.Sleep(time.Second * time.Duration(g.Conf().Strategy.UpdateDuration))
	}
}

// GetOldestTms to analysis the oldes tms
func GetOldestTms(filepath string) (int64, bool) {
	JobManagerLock.RLock()
	defer JobManagerLock.RUnlock()
	job, ok := JobManager[filepath]
	if !ok {
		return 0, false
	}

	tms, allFree := job.w.GetOldestTms()
	nowTms := time.Now().Unix()
	//如果worker全都空闲，且当前时间戳已经领先1min
	//则将标记的时间戳设置为当前时间戳-1min
	if allFree && nowTms-tms > 60 {
		tms = nowTms - 60
	}
	return tms, true
}

//添加任务到管理map( managerjob managerconfig) 启动reader和worker
func createJob(config *ConfigInfo, cache chan string, st *scheme.Strategy) error {
	if _, ok := JobManager[config.FilePath]; ok {
		if _, ok := ManagerConfig[config.ID]; !ok {
			ManagerConfig[config.ID] = config
		}
		return nil
	}

	ManagerConfig[config.ID] = config
	//启动reader
	r, err := reader.NewReader(config.FilePath, cache)
	if err != nil {
		return err
	}
	dlog.Infof("Add Reader : [%s]", config.FilePath)
	//启动worker
	w := NewWorkerGroup(config.FilePath, cache, st)
	JobManager[config.FilePath] = &Job{
		r: r,
		w: w,
	}
	w.Start()
	//启动reader
	go r.Start()

	dlog.Infof("Create job success [filePath:%s][sid:%d]", config.FilePath, config.ID)
	return nil
}

//先stop worker reader再从管理map中删除
func deleteJob(config *ConfigInfo) {
	//删除jobs
	tag := 0
	for _, cg := range ManagerConfig {
		if config.FilePath == cg.FilePath {
			tag++
		}
	}
	if tag <= 1 {
		dlog.Infof("Del Reader : [%s]", config.FilePath)
		if job, ok := JobManager[config.FilePath]; ok {
			job.w.Stop() //先stop worker
			job.r.Stop()
			delete(JobManager, config.FilePath)
		}
	}
	dlog.Infof("Stop reader & worker success [filePath:%s][sid:%d]", config.FilePath, config.ID)

	//删除config
	if _, ok := ManagerConfig[config.ID]; ok {
		delete(ManagerConfig, config.ID)
	}
}
