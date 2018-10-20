package reader

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/proc/metric"

	"github.com/hpcloud/tail"
	"github.com/open-falcon/falcon-plus/modules/agent/collector"
	"regexp"
)

// Reader to read file
type Reader struct {
	FilePath    string //配置的路径 正则路径
	t           *tail.Tail
	Stream      chan string
	CurrentPath string //当前的路径
	Close       chan struct{}
}

// NewReader to create a reader
func NewReader(filepath string, stream chan string, prefix string) (*Reader, error) {
	r := &Reader{
		FilePath: filepath,
		Stream:   stream,
		Close:    make(chan struct{}),
	}

	reg, err := regexp.Compile(prefix)
	if err != nil {
		return nil, err
	}

	path := GetCurrentPath(filepath)
	id := generateId(path)

	go func() {
		for {
			path = GetCurrentPath(filepath)
			collector.Read(id, path, *reg, stream)
			time.Sleep(60 * time.Second)
		}
	}()

	return r, err
}

func generateId(str string) uint64 {
	var id = uint64(0)
	for _, v := range []byte(str) {
		id += uint64(v)
	}
	return id
}

// StartRead to start to read
func (r *Reader) StartRead() {
	var readCnt, readSwp int64
	var dropCnt, dropSwp int64

	go func() {
		for {
			// 十秒钟统计一次
			select {
			case <-r.Close:
				return
			case <-time.After(time.Second * 10):
			}
			// 统计时间戳可以不准，但是不能漏
			a := readCnt
			b := dropCnt
			metric.MetricReadLine(r.FilePath, a-readSwp)
			metric.MetricDropLine(r.FilePath, b-dropSwp)
			readSwp = a
			dropSwp = b
		}
	}()
}

// StopRead to stop a read instance
func (r *Reader) StopRead() error {
	return r.t.Stop()
}

// Stop to stop a reader
func (r *Reader) Stop() {
	r.StopRead()
	close(r.Close)

}

// Start a reader
func (r *Reader) Start() {
	go r.StartRead()
	for {
		select {
		case <-time.After(time.Second):
			r.check()
		case <-r.Close:
			close(r.Stream)
			return
		}
	}

}

func (r *Reader) check() {
}
