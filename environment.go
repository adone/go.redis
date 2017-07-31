package redis

import (
	"fmt"
	"net/url"
	"os"
	"time"
)

// ENV возвращает конфигурацию подключения
func ENV(prefix string) Configuration {
	config := Configuration{
		Timeout:        CommonTimeout(prefix),
		ConnectTimeout: ConnectionTimeout(prefix),
		ReadTimeout:    ReadTimeout(prefix),
		WriteTimeout:   WriteTimeout(prefix),
	}

	config.fromENV(prefix)

	return config
}

// Составляем url подключения к редису
func (config *Configuration) fromENV(prefix string) {
	address, database, user, password := GetENVData(prefix)

	config.Database = database
	config.User = user
	config.Password = password
	config.url = address

	sentinelURLs, masterName := ParseSentinelENVs(prefix)
	if len(sentinelURLs) != 0 {
		config.Sentinel = NewSentinel(config, sentinelURLs, masterName)
	}

}

func GetENVData(prefix string) (address string, database string, user string, password string) {
	address, database, user, password, host, port := GetENVs(prefix)
	if address == "" {
		hostPort := fmt.Sprintf("%s:%s", host, port)
		address = CompileURL(hostPort, database, user, password)
	}
	return address, database, user, password
}

func GetENVs(prefix string) (address string, database string, user string, password string, host string, port string) {
	address = os.Getenv(fmt.Sprintf("%s_REDIS_URL", prefix))
	database = os.Getenv(fmt.Sprintf("%s_REDIS_DATABASE", prefix))
	user = os.Getenv(fmt.Sprintf("%s_REDIS_USERNAME", prefix))
	password = os.Getenv(fmt.Sprintf("%s_REDIS_PASSWORD", prefix))

	host = os.Getenv(fmt.Sprintf("%s_REDIS_SERVICE_HOST", prefix))
	if host == "" {
		host = "localhost"
	}

	port = os.Getenv(fmt.Sprintf("%s_REDIS_SERVICE_PORT", prefix))
	if port == "" {
		port = "6379"
	}

	return address, database, user, password, host, port
}

func CompileURL(hostPort, database, user, password string) string {
	address := url.URL{
		Scheme: "redis",
		Path:   database,
	}
	address.Host = hostPort

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
