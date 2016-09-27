package redis

import (
	"fmt"
	"net/url"
	"os"
	"time"
)

// ENV возвращает конфигурацию подключения
func ENV(prefix string) Configuration {
	return Configuration{
		URL:            URL(prefix),
		Timeout:        CommonTimeout(prefix),
		ConnectTimeout: ConnectionTimeout(prefix),
		ReadTimeout:    ReadTimeout(prefix),
		WriteTimeout:   WriteTimeout(prefix),
	}
}

// URL урл подключения к редису
func URL(prefix string) string {
	if address := os.Getenv(fmt.Sprintf("%s_REDIS_URL", prefix)); address != "" {
		return address
	}

	address := url.URL{
		Scheme: "redis",
		Path:   os.Getenv(fmt.Sprintf("%s_REDIS_DATABASE", prefix)),
	}

	host := os.Getenv(fmt.Sprintf("%s_REDIS_SERVICE_HOST", prefix))
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv(fmt.Sprintf("%s_REDIS_SERVICE_PORT", prefix))
	if port == "" {
		port = "6379"
	}

	address.Host = fmt.Sprintf("%s:%s", host, port)

	user := os.Getenv(fmt.Sprintf("%s_REDIS_USERNAME", prefix))
	password := os.Getenv(fmt.Sprintf("%s_REDIS_PASSWORD", prefix))

	if user != "" && password != "" {
		address.User = url.UserPassword(user, password)
	}

	return address.String()
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
