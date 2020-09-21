package downloader

import (
	"strconv"
	"strings"
)

func decodeContentRange(str string) (unit string, s, e, l int64) {
	s, e, l = -1, -1, -1
	parts := strings.Fields(str)
	if len(parts) != 2 {
		return
	}
	unit = parts[0]
	parts = strings.Split(parts[1], "/")
	if len(parts) != 2 {
		return
	}
	if v, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
		l = int64(v)
	}

	parts = strings.Split(parts[0], "-")
	if len(parts) != 2 {
		return
	}
	if v, err := strconv.ParseUint(parts[0], 10, 64); err == nil {
		s = int64(v)
	}
	if v, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
		e = int64(v)
	}
	return
}
