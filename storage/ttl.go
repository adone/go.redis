package storage

import (
	"time"
)

// TTL время жизни ключа
type TTL struct {
	Key   string
	Value interface{}
}

// Seconds время жизни в секундах
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
