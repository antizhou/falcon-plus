package reader

import (
	"encoding/json"
	"sync"
	"time"
	"syscall"
	log "github.com/open-falcon/falcon-plus/logger"
)

const (
	PersistenceFilename = "collector.json"
)

var (
	PersistenceMutex  = sync.RWMutex{}
	PersistenceLoaded = false
	PersistenceData   []PersistenceRow
)

type PersistenceRow struct {
	ID        uint64 `json:"id"`
	Path      string `json:"path"`
	Offset    int64  `json:"offset"`
	Device    uint64 `json:"device"`
	Inode     uint64 `json:"inode"`
	Timestamp int64  `json:"timestamp"`
}

func NewOffset(id uint64, filePath string) (r *PersistenceRow) {
	PersistenceMutex.Lock()
	defer PersistenceMutex.Unlock()

	var stat syscall.Stat_t
	if err := syscall.Stat(filePath, &stat); err != nil {
		log.Errorf("stat file error", err)
		return
	}

	row := PersistenceRow{
		ID:        id,
		Path:      filePath,
		Offset:    0,
		Device:    uint64(stat.Dev),
		Inode:     stat.Ino,
		Timestamp: time.Now().Unix(),
	}
	PersistenceData = append(PersistenceData, row)

	r = &row
	return
}

func loadOffset(id uint64, path string) (r *PersistenceRow) {
	err := loadPersistence()

	if err != nil {
		return
	}

	PersistenceMutex.RLock()
	defer PersistenceMutex.RUnlock()
	for _, row := range PersistenceData {
		if row.ID == id && row.Path == path {
			r = &row
			return
		}
	}

	return
}

func SaveOffset(id uint64, path string, offset int64, device, inode uint64) {
	idx := -1

	PersistenceMutex.RLock()
	for i, row := range PersistenceData {
		if row.ID == id && row.Path == path {
			idx = i
		}
	}
	PersistenceMutex.RUnlock()

	PersistenceMutex.Lock()
	if idx != -1 {
		PersistenceData[idx].Offset = offset
		PersistenceData[idx].Device = device
		PersistenceData[idx].Inode = inode
		PersistenceData[idx].Timestamp = time.Now().Unix()
	} else {
		PersistenceData = append(PersistenceData, PersistenceRow{
			ID:        id,
			Path:      path,
			Offset:    offset,
			Device:    device,
			Inode:     inode,
			Timestamp: time.Now().Unix(),
		})
	}
	PersistenceMutex.Unlock()

	savePersistence()

	return
}

func loadPersistence() error {
	PersistenceMutex.Lock()
	defer PersistenceMutex.Unlock()

	if PersistenceLoaded {
		return nil
	}

	PersistenceLoaded = true
	data, err := Load(PersistenceFilename)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &PersistenceData)

	if err != nil {
		return err
	}

	return nil
}

func savePersistence() error {
	PersistenceMutex.Lock()
	defer PersistenceMutex.Unlock()

	data, err := json.Marshal(&PersistenceData)

	if err != nil {
		return err
	}

	err = Save(PersistenceFilename, data)

	if err != nil {
		log.Error("log collector can not save offset, err: ", err)
	}

	return err
}

func ClearPersistenceData() {
	PersistenceMutex.Lock()
	defer PersistenceMutex.Unlock()

	PersistenceData = make([]PersistenceRow, 0)
}
