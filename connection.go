package redis

import (
	"github.com/garyburd/redigo/redis"
)

// New создание подключения к редису по переданным настройкам
func New(config Configuration) (redis.Conn, error) {
	return redis.DialURL(config.URL,
		redis.DialConnectTimeout(config.GetConnectTimeout()),
		redis.DialReadTimeout(config.GetReadTimeout()),
		redis.DialWriteTimeout(config.GetWriteTimeout()),
		redis.DialDatabase(config.Database),
	)
}

// Connect подключение к редису
func Connect(configuration Configuration) func() (redis.Conn, error) {
	return func() (redis.Conn, error) { return New(configuration) }
}
