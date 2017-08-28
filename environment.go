package redis

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// ENV return configuration from env variables
func ENV(prefix string) *Configuration {
	config := &Configuration{
		Timeout:           CommonTimeout(prefix),
		ConnectTimeout:    ConnectionTimeout(prefix),
		ReadTimeout:       ReadTimeout(prefix),
		WriteTimeout:      WriteTimeout(prefix),
		MasterName:        MasterName(prefix),
		SentinelAddresses: SentinelAddresses(prefix),
		Password:          os.Getenv(fmt.Sprintf("%s_REDIS_PASSWORD", prefix)),
		Database:          Database(prefix),
	}

	// TODO: move it to Dialer
	if len(config.SentinelAddresses) > 0 {
		config.Sentinel = NewSentinel(config)
		return config
	}

	config.address = os.Getenv(fmt.Sprintf("%s_REDIS_ADDRESS", prefix))
	if config.address == "" {
		host := os.Getenv(fmt.Sprintf("%s_REDIS_SERVICE_HOST", prefix))
		if host == "" {
			host = "localhost"
		}

		port := os.Getenv(fmt.Sprintf("%s_REDIS_SERVICE_PORT", prefix))
		if port == "" {
			port = "6379"
		}

		config.address = net.JoinHostPort(host, port)
	}

	return config
}

func SentinelAddresses(prefix string) []string {
	addresses := os.Getenv(fmt.Sprintf("%s_REDIS_SENTINEL_ADDRESSES", prefix))

	if addresses == "" {
		return nil
	}

	return strings.Split(addresses, ",")
}

func MasterName(prefix string) string {
	master := os.Getenv(fmt.Sprintf("%s_REDIS_SENTINEL_MASTER_NAME", prefix))
	if master == "" {
		master = os.Getenv("REDIS_SENTINEL_MASTER_NAME")
	}

	if master == "" {
		return "mymaster"
	}

	return master
}

func Database(prefix string) int {
	if database, err := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_REDIS_DATABASE", prefix))); err == nil {
		return database
	}

	return 0
}

// CommonTimeout общий таймаут
func CommonTimeout(prefix string) time.Duration {
	value := os.Getenv(fmt.Sprintf("%s_REDIS_TIMEOUT", prefix))
	if value == "" {
		value = os.Getenv("REDIS_TIMEOUT")
	}

	if timeout, err := time.ParseDuration(value); err == nil {
		return timeout
	}

	return 0
}

// ConnectionTimeout таймаут на подключение
func ConnectionTimeout(prefix string) time.Duration {
	value := os.Getenv(fmt.Sprintf("%s_REDIS_CONNECT_TIMEOUT", prefix))
	if value == "" {
		value = os.Getenv("REDIS_CONNECT_TIMEOUT")
	}

	if timeout, err := time.ParseDuration(value); err == nil {
		return timeout
	}

	return 0
}

// ReadTimeout таймаут на чтение
func ReadTimeout(prefix string) time.Duration {
	value := os.Getenv(fmt.Sprintf("%s_REDIS_READ_TIMEOUT", prefix))
	if value == "" {
		value = os.Getenv("REDIS_READ_TIMEOUT")
	}

	if timeout, err := time.ParseDuration(value); err == nil {
		return timeout
	}

	return 0
}

// WriteTimeout таймаут на запись
func WriteTimeout(prefix string) time.Duration {
	value := os.Getenv(fmt.Sprintf("%s_REDIS_WRITE_TIMEOUT", prefix))
	if value == "" {
		value = os.Getenv("REDIS_WRITE_TIMEOUT")
	}

	if timeout, err := time.ParseDuration(value); err == nil {
		return timeout
	}

	return 0
}
