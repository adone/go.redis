package storage

import (
	"github.com/garyburd/redigo/redis"
)

type Configuration struct {
	KeyTTL    interface{} // Common key time-to-live, if set affects every key used in storage
	Namespace string

	Pool       *redis.Pool
	Connection redis.Conn
}
