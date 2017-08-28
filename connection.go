package redis

import (
	"github.com/FZambia/go-sentinel"
	"github.com/garyburd/redigo/redis"
)

// NewSentinel creates new Sentinel connection
func NewSentinel(config *Configuration) *sentinel.Sentinel {
	dialer := NewDialer(config)

	return &sentinel.Sentinel{
		Addrs:      config.SentinelAddresses,
		MasterName: config.MasterName,
		Dial:       dialer.Dial,
	}
}

// New creates new redis connection
func New(config *Configuration) (redis.Conn, error) {
	address, err := config.Address()
	if err != nil {
		return nil, err
	}

	return NewDialer(config).Dial(address)
}

func NewDialer(config *Configuration) *Dialer {
	options := make([]redis.DialOption, 0, 5)

	options = append(options,
		redis.DialConnectTimeout(config.GetConnectTimeout()),
		redis.DialReadTimeout(config.GetReadTimeout()),
		redis.DialWriteTimeout(config.GetWriteTimeout()),
		redis.DialDatabase(config.Database),
	)

	if config.Password != "" {
		options = append(options, redis.DialPassword(config.Password))
	}

	return &Dialer{options}
}

type Dialer struct {
	options []redis.DialOption
}

func (dialer Dialer) Dial(address string) (redis.Conn, error) {
	return redis.Dial("tcp", address, dialer.options...)
}

func Connect(configuration *Configuration) func() (redis.Conn, error) {
	dialer := NewDialer(configuration)

	return func() (redis.Conn, error) {
		address, err := configuration.Address()
		if err != nil {
			return nil, err
		}

		return dialer.Dial(address)
	}
}
