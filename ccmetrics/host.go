package ccmetrics

import (
	"regexp"
	"strings"
)

type host string

func (h host) String() string {
	return string(h)
}

func (h host) Shard() string {
	host := h.String()
	parts := strings.Split(host, ".")
	if len(parts) == 2 {
		s := regexp.MustCompile(`[0-9]+`).FindString(parts[0])
		if s != "" {
			return s
		}
	}
	return host
}
