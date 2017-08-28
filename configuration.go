package redis

import (
	"github.com/FZambia/go-sentinel"
	"time"
)

// Configuration
type Configuration struct {
	address           string
	MasterName        string
	SentinelAddresses []string
	Sentinel          *sentinel.Sentinel
	Timeout           time.Duration
	ConnectTimeout    time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	Database          int
	Password          string
}

// GetConnectTimeout returns connection timeout
func (config Configuration) GetConnectTimeout() time.Duration {
	if config.ConnectTimeout == 0 {
		return config.Timeout
	}

	return config.ConnectTimeout
}

// GetReadTimeout returns response read timeout
func (config Configuration) GetReadTimeout() time.Duration {
	if config.ReadTimeout == 0 {
		return config.Timeout
	}

	return config.ReadTimeout
}

// GetWriteTimeout returns request write timeout
func (config Configuration) GetWriteTimeout() time.Duration {
	if config.WriteTimeout == 0 {
		return config.Timeout
	}

	return config.WriteTimeout
}

// Address return redis addess
func (config Configuration) Address() (string, error) {
	if config.Sentinel == nil {
		return config.address, nil
	}

	return config.Sentinel.MasterAddr()
}
