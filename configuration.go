package redis

import (
	"time"
)

// Configuration настройки подключения к редису
type Configuration struct {
	URL            string        // Адрес хоста редиса
	Timeout        time.Duration // Общий таймаут
	ConnectTimeout time.Duration // Таймаут на подключение
	ReadTimeout    time.Duration // Таймаут на чтение
	WriteTimeout   time.Duration // Таймаут на запись
	Database       int           // Номер базы данных
}

// GetConnectTimeout получение размера таймаута на подключение
func (config Configuration) GetConnectTimeout() time.Duration {
	if config.ConnectTimeout == 0 {
		return config.Timeout
	}

	return config.ConnectTimeout
}

// GetConnectTimeout получение размера таймаута на чтение
func (config Configuration) GetReadTimeout() time.Duration {
	if config.ReadTimeout == 0 {
		return config.Timeout
	}

	return config.ReadTimeout
}

// GetConnectTimeout получение размера таймаута на запись
func (config Configuration) GetWriteTimeout() time.Duration {
	if config.WriteTimeout == 0 {
		return config.Timeout
	}

	return config.WriteTimeout
}
