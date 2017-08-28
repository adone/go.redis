package pool

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// New creates new Redis connection pool
func New(config Configuration,
	dial func() (redis.Conn, error),
	check func(redis.Conn, time.Time) error,
) *redis.Pool {
	return &redis.Pool{
		Wait:         config.WaitConnection,
		MaxIdle:      config.MaxIdleConnectionCount,
		MaxActive:    config.MaxActiveConnectionCount,
		IdleTimeout:  config.IdleConnectionTimeout,
		Dial:         dial,
		TestOnBorrow: check,
	}
}

func Check(configuration Configuration) func(redis.Conn, time.Time) error {
	return func(connection redis.Conn, previous time.Time) error {
		if configuration.CheckConnectionFrequency == 0 {
			return nil
		}

		if time.Since(previous) < configuration.CheckConnectionFrequency {
			return nil
		}

		_, err := connection.Do("PING")
		return err
	}
}
