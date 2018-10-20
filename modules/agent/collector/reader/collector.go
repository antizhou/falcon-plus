package reader

import (
	"os"
	"regexp"
	"time"

	log "github.com/open-falcon/falcon-plus/logger"
)

func Read(id uint64, filePath string, prefix regexp.Regexp, stream chan string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		log.Errorf("[agent.collector.log.Read] file stat error: ", err)
		return err
	}

	r := loadOffset(id, filePath)
	if r == nil {
		r = NewOffset(id, filePath)
	}

	size := stat.Size()
	left := size - r.Offset

	if left == 0 {
		time.Sleep(10 * time.Second)
		return nil
	}

	log.Infof("[agent.collector.log.Read] the file %s has %d bytes left", r.Path, left)

	for {
		if r.Offset >= size {
			break
		}

		o, err := ReadLines(r, file, size, prefix, stream)
		if o == 0 {
			time.Sleep(10 * time.Second)
			break
		}
		if err != nil {
			log.Error("can not read data on ", r.Path, " err: ", err)
			time.Sleep(10 * time.Second)
			break
		}
	}

	return nil
}

func ReadLines(r *PersistenceRow, file *os.File, size int64, prefix regexp.Regexp, stream chan string) (int, error) {
	currentOffset := 0

	buf := make([]byte, fileBufLen)

	left := int(size - r.Offset)

	currentLen := fileBufLen

	if currentLen > left {
		currentLen = left
	}

	log.Infof("[agent.collector.log.ReadLines] collector reading %d bytes on %s", currentLen, r.Path)

	file.Seek(r.Offset, os.SEEK_SET)
	_, err := file.Read(buf[:currentLen])

	if err != nil {
		log.Errorf("[agent.collector.log.ReadLines] collector can not read file %s, err: %s", r.Path, err.Error())
		return currentOffset, err
	}

	if err != nil {
		log.Error("[agent.collector.log.ReadLines] can not parse prefix regexp, err: ", err)
		return currentOffset, err
	}

	addOffset := 0

	for {
		if currentOffset >= currentLen {
			break
		}

		line, o := readLine(buf[currentOffset:currentLen], &prefix)
		stream <- string(line)

		if o == 0 {
			break
		}

		currentOffset += o
		addOffset += o
	}

	if addOffset > 0 {
		r.Offset += int64(addOffset)
		SaveOffset(r.ID, r.Path, r.Offset, r.Device, r.Inode)
		addOffset = 0
	}

	return currentOffset, nil
}
