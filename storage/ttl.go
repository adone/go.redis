package storage

import (
	"time"
)

type TTL struct {
	Key   string
	Value interface{}
}

// Seconds returns converted time from initial value in seconds
func (ttl TTL) Seconds() int {
	switch value := ttl.Value.(type) {
	case int:
		return value
	case time.Duration:
		return int(value / time.Second)
	case func(string) int:
		return value(ttl.Key)
	default:
		return 0
	}
}
