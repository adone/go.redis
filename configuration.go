package redis

import (
	"github.com/FZambia/go-sentinel"
	"time"
)

// Configuration настройки подключения к редису
type Configuration struct {
	url            string             // Адрес хоста редиса
	Sentinel       *sentinel.Sentinel // Подключение к Sentinel
	Timeout        time.Duration      // Общий таймаут
	ConnectTimeout time.Duration      // Таймаут на подключение
	ReadTimeout    time.Duration      // Таймаут на чтение
	WriteTimeout   time.Duration      // Таймаут на запись
	Database       string             // Номер базы данных
	User           string             // Пользователь базы данных
	Password       string             // Пароль пользователя базы данных
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

// URL выдает URL для соединения с Redis
// Либо прямую ссылку, либо через запрос к Sentinel
func (config Configuration) URL() (string, error) {
	url := config.url

	if config.Sentinel != nil {
		addr, err := config.Sentinel.MasterAddr()
		if err != nil {
			return "", err
		}
		url = CompileURL(addr, config.Database, config.User, config.Password)
	}

	return url, nil
}
