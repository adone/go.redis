package storage

import (
	"github.com/garyburd/redigo/redis"
)

// Configuration настройки стораджа кеша
type Configuration struct {
	KeyTTL    interface{} // Время жизни ключа
	Namespace string

	Pool       *redis.Pool
	Connection redis.Conn
}
