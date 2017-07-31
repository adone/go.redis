package redis

import (
	"fmt"
	"github.com/FZambia/go-sentinel"
	"github.com/garyburd/redigo/redis"
	"os"
	"strings"
)

// ParseSentinelENVs парсит адреса Sentinel из ENV
// Читает prefix_SENTINEL_MASTER_NAME, либо задает его значением по умолчанию "mymaster"
func ParseSentinelENVs(prefix string) ([]string, string) {
	masterName := os.Getenv(fmt.Sprintf("%s_SENTINEL_MASTER_NAME", prefix))
	if masterName == "" {
		masterName = "mymaster"
	}
	if address := os.Getenv(fmt.Sprintf("%s_SENTINEL_HOSTS_PORTS", prefix)); address != "" {
		return strings.Split(address, ","), masterName
	}
	return nil, ""
}

// NewSentinel создает новый объект подключения к Sentinel
func NewSentinel(config *Configuration, urls []string, masterName string) *sentinel.Sentinel {
	return &sentinel.Sentinel{
		Addrs:      urls,
		MasterName: masterName,
		Dial: func(addr string) (redis.Conn, error) {
			return redis.Dial("tcp", addr,
				redis.DialConnectTimeout(config.GetConnectTimeout()),
				redis.DialReadTimeout(config.GetReadTimeout()),
				redis.DialWriteTimeout(config.GetWriteTimeout()),
			)
		},
	}
}
