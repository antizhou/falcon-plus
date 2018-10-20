package reader

import (
	"bytes"
	"regexp"
)

const (
	// 10M
	fileBufLen = 1024 * 1024 * 10
)

var (
	lineSep    = []byte("\n")
	lineSepLen = len(lineSep)
)

func readLine(content []byte, prefix *regexp.Regexp) ([]byte, int) {
	e := 0

	for {
		if len(content) <= e {
			break
		}

		i := bytes.Index(content[e:], lineSep)

		if i == -1 {
			break
		}

		if prefix == nil {
			e = e + i + lineSepLen
			break
		}

		if e == 0 {
			e += i + lineSepLen
			continue
		}

		t := e + i + lineSepLen

		if prefix.Match(content[e:t]) {
			break
		}

		e = t
	}

	if e == 0 {
		return nil, 0
	}

	data := content[:e]
	length := len(data)

	return data, length
}
