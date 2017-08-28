package storage

import (
	"github.com/garyburd/redigo/redis"
)

const (
	START = "0"
	SCAN  = "SCAN"
	SSCAN = "SSCAN"
)

type Option func(*Iterator)

func NewIterator(options ...Option) *Iterator {
	iterator := new(Iterator)
	iterator.command = SCAN
	iterator.cursor = START

	for _, option := range options {
		option(iterator)
	}

	return iterator
}

func WithStorage(storage *Client) Option {
	return func(iterator *Iterator) {
		iterator.storage = storage
	}
}

func WithCursor(cursor string) Option {
	return func(iterator *Iterator) {
		iterator.cursor = cursor
	}
}

func ForSet(key string) Option {
	return func(iterator *Iterator) {
		iterator.command = SSCAN
		iterator.key = key
	}
}

func WithTemplate(template string) Option {
	return func(iterator *Iterator) {
		iterator.template = template
	}
}

func WithBatchSize(batchSize int) Option {
	return func(iterator *Iterator) {
		iterator.batchSize = batchSize
	}
}

type Iterator struct {
	command   string
	key       string
	cursor    string
	storage   *Client
	template  string
	batchSize int
}

// All iterates all keys with provided function
func (iterator *Iterator) All(yield func([]interface{})) error {
	if iterator.batchSize == 0 {
		return nil
	}

	connection := iterator.storage.checkout()
	defer iterator.storage.release(connection)

	for data, err := iterator.next(connection); ; data, err = iterator.next(connection) {
		if err != nil {
			return err
		}

		yield(iterator.handle(data))

		if iterator.cursor == START {
			break
		}
	}

	return nil
}

func (iterator *Iterator) Next() ([]interface{}, error) {
	connection := iterator.storage.checkout()
	defer iterator.storage.release(connection)

	data, err := iterator.next(connection)
	if err != nil {
		return nil, err
	}

	return iterator.handle(data), nil
}

func (iterator *Iterator) handle(data interface{}) (result []interface{}) {
	// response format - [cursor,[value,value,...]]
	if results, ok := data.([]interface{}); ok && len(results) == 2 {
		if values, ok := results[1].([]interface{}); ok {
			result = values
		}

		if next, ok := results[0].([]byte); ok {
			iterator.cursor = string(next)
		}
	}

	return
}

func (iterator *Iterator) next(connection redis.Conn) (interface{}, error) {
	args := make([]interface{}, 0, 6)

	if iterator.key != "" {
		args = append(args, iterator.key)
	}

	args = append(args, iterator.cursor)

	if iterator.template != "" {
		args = append(args, "MATCH", iterator.template)
	}

	args = append(args, "COUNT", iterator.batchSize)

	return connection.Do(iterator.command, args...)
}
