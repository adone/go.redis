package redis

import (
	"github.com/garyburd/redigo/redis"
)

// New создание подключения к редису по переданным настройкам
func New(config Configuration) (redis.Conn, error) {
	url, err := config.URL()
	if err != nil {
		return nil, err
	}
	return redis.DialURL(url,
		redis.DialConnectTimeout(config.GetConnectTimeout()),
		redis.DialReadTimeout(config.GetReadTimeout()),
		redis.DialWriteTimeout(config.GetWriteTimeout()),
	)
}

// Connect подключение к редису
func Connect(configuration Configuration) func() (redis.Conn, error) {
	return func() (redis.Conn, error) { return New(configuration) }
}
