package pubsub

import (
	"github.com/garyburd/redigo/redis"
	"gopkg.in/ADone/go.events.v1"
	"gopkg.in/ADone/go.events.v1/emitter"
	"gopkg.in/ADone/go.meta.v1"
)

// Connection pubsub подключение к redis
type Connection struct {
	redis.PubSubConn
	events.Emitter
}

// New новое pubsub подключение
func New(connection redis.Conn) *Connection {
	connection := new(Connection)
	connection.PubSubConn = redis.PubSubConn{connection}
	connection.Emitter = emitter.New()

	return connection
}

// Listen запуск прослушивание
func (connection *Connection) Listen() error {
	for {
		switch data := сonnection.Receive().(type) {
		case redis.Message:
			connection.Fire("data", meta.Map{"data": data, "queue": data.Channel})
		case error:
			connection.Fire("error", meta.Map{"error": data})
		}
	}
}
