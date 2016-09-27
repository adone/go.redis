package pool

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	// DefaultRedisPoolSize Максимальное количество свободных соединений в пулле
	DefaultRedisPoolSize = 8
)

// ENV возвращает конфигурацию пулла
func ENV(prefix string) Configuration {
	return Configuration{
		WaitConnection:           true,
		MaxIdleConnectionCount:   MaxIdleCount(prefix),
		MaxActiveConnectionCount: MaxActiveCount(prefix),
		IdleConnectionTimeout:    IdleTimeout(prefix),
		CheckConnectionFrequency: CheckFrequency(prefix),
	}
}

// MaxActiveCount максимальное количество активных соединений
func MaxActiveCount(prefix string) int {
	value := os.Getenv(fmt.Sprintf("%s_REDIS_ACTIVE_POOL_SIZE", prefix))
	if value == "" {
		value = os.Getenv("REDIS_ACTIVE_POOL_SIZE")
	}

	if size, err := strconv.Atoi(value); err == nil {
		return size
	}

	return 0
}

// MaxIdleCount максимальное количестов соединений в режиме ожидания
func MaxIdleCount(prefix string) int {
	value := os.Getenv(fmt.Sprintf("%s_REDIS_IDLE_POOL_SIZE", prefix))
	if value == "" {
		value = os.Getenv("REDIS_IDLE_POOL_SIZE")
	}

	if value == "" {
		value = os.Getenv("REDIS_POOL_SIZE")
	}

	if size, err := strconv.Atoi(value); err == nil {
		return size
	}

	return DefaultRedisPoolSize
}

// IdleTimeout время жизни соединения в пулле
func IdleTimeout(prefix string) time.Duration {
	value := os.Getenv(fmt.Sprintf("%s_REDIS_POOL_TIMEOUT", prefix))
	if value == "" {
		value = os.Getenv("REDIS_POOL_IDLE_TIMEOUT")
	}

	if value == "" {
		value = os.Getenv("REDIS_POOL_TIMEOUT")
	}

	if timeout, err := time.ParseDuration(value); err == nil {
		return timeout
	}

	return 0
}

// CheckFrequency частота проверки соединения
func CheckFrequency(prefix string) time.Duration {
	value := os.Getenv(fmt.Sprintf("%s_REDIS_POOL_CHECK_TIMEOUT", prefix))
	if value == "" {
		value = os.Getenv("REDIS_POOL_CHECK_TIMEOUT")
	}

	if timeout, err := time.ParseDuration(value); err == nil {
		return timeout
	}

	return 0
}
