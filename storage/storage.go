package storage

import (
	"sync"

	"github.com/garyburd/redigo/redis"
)

// New creates new Redis client
func New(config Configuration) *Client {
	storage := &Client{
		guard:     new(sync.Mutex),
		KeyTTL:    config.KeyTTL,
		Namespace: config.Namespace,
	}

	if config.Pool != nil {
		storage.pool = config.Pool
		return storage
	}

	if config.Connection != nil {
		storage.connection = config.Connection
		return storage
	}

	panic("redis storage: no connection provided")

	return nil
}

type Client struct {
	KeyTTL    interface{} // Common key time-to-live, if set affects every key used in storage
	Namespace string

	pool       *redis.Pool
	guard      *sync.Mutex
	connection redis.Conn
}

func (storage *Client) checkout() redis.Conn {
	if storage.pool != nil {
		return storage.pool.Get()
	}

	storage.guard.Lock()
	return storage.connection
}

func (storage *Client) release(connection redis.Conn) {
	if storage.pool != nil {
		connection.Close()
		return
	}

	storage.guard.Unlock()
}

// Expire see EXPIRE
func (storage *Client) Expire(key string, ttl interface{}) error {
	connection := storage.checkout()
	defer storage.release(connection)

	_, err := connection.Do("EXPIRE", key, TTL{key, ttl}.Seconds())
	return err
}

// Set see SET
func (storage *Client) Set(key string, value []byte) error {
	setter := Setter{
		Storage: storage,
		TTL:     storage.KeyTTL,
		Key:     key,
		Value:   value,
	}

	return setter.Call()
}

// Increment see INCREMENT
func (storage *Client) Increment(key string, delta int) (int, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	return redis.Int(connection.Do("INCRBY", key, delta))
}

// Get see GET
func (storage *Client) Get(key string) ([]byte, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	data, err := redis.Bytes(connection.Do("GET", key))
	if err == redis.ErrNil {
		return []byte{}, nil
	}

	return data, err
}

// MultiGet see MGET
func (storage *Client) MultiGet(keys ...string) ([][]byte, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	args := make([]interface{}, len(keys))
	for index, key := range keys {
		args[index] = key
	}

	data, err := redis.ByteSlices(connection.Do("MGET", args...))
	if err == redis.ErrNil {
		return [][]byte{}, nil
	}

	return data, err
}

// Publish see PUBLISH
func (storage *Client) Publish(key string, value []byte) error {
	connection := storage.checkout()
	defer storage.release(connection)

	_, err := connection.Do("PUBLISH", key, value)

	return err
}

// Keys see SCAN, it does not use KEYS because it recommended by Redis team https://redis.io/commands/keys
func (storage *Client) Keys(template string) ([]string, error) {
	iterator := NewIterator(WithStorage(storage), WithTemplate(template), WithBatchSize(32))

	var keys []string

	return keys, iterator.All(func(values []interface{}) {
		if found, err := redis.Strings(values, nil); err == nil {
			keys = append(keys, found...)
		}
	})
}

// SetField see HSET
func (storage *Client) SetField(key, field string, value []byte) error {
	connection := storage.checkout()
	defer storage.release(connection)

	_, err := connection.Do("HSET", key, field, value)

	return err
}

// GetField see HGET
func (storage *Client) GetField(key, field string) ([]byte, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	data, err := redis.Bytes(connection.Do("HGET", key, field))

	if err == redis.ErrNil {
		return []byte{}, nil
	}

	return data, err
}

// SetFields see HMSET
func (storage *Client) SetFields(key string, hash map[string]interface{}) error {
	if len(hash) == 0 {
		return nil
	}

	connection := storage.checkout()
	defer storage.release(connection)

	args := make([]interface{}, 2*len(hash)+1)
	args[0] = key
	index := 1

	for field, value := range hash {
		args[index] = field
		args[index+1] = value
		index += 2
	}

	_, err := connection.Do("HMSET", args...)

	return err
}

// GetFields see HMGET
func (storage *Client) GetFields(keyAndFields ...string) (map[string][]byte, error) {
	if len(keyAndFields) <= 1 {
		return nil, nil
	}

	connection := storage.checkout()
	defer storage.release(connection)

	args := make([]interface{}, len(keyAndFields))
	for index, value := range keyAndFields {
		args[index] = value
	}

	data, err := redis.ByteSlices(connection.Do("HMGET", args...))

	hash := make(map[string][]byte)
	for index, value := range data {
		hash[keyAndFields[index+1]] = value
	}

	return hash, err
}

// IncrementField see HINCRBY
func (storage *Client) IncrementField(key, field string, delta int) (int, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	return redis.Int(connection.Do("HINCRBY", key, field, delta))
}

// FieldExist see HEXISTS
func (storage *Client) FieldExist(key, field string) (bool, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	return redis.Bool(connection.Do("HEXISTS", key, field))
}

// GetValues see HVALS
func (storage *Client) GetValues(key string) ([][]byte, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	data, err := redis.ByteSlices(connection.Do("HVALS", key))

	if err == redis.ErrNil {
		return [][]byte{}, nil
	}

	return data, err
}

// RemoveFields see HDEL
func (storage *Client) RemoveFields(keyAndFields ...string) error {
	if len(keyAndFields) <= 1 {
		return nil
	}

	connection := storage.checkout()
	defer storage.release(connection)

	args := make([]interface{}, len(keyAndFields))
	for index, keyOrField := range keyAndFields {
		args[index] = keyOrField
	}

	_, err := connection.Do("HDEL", args...)

	return err
}

// Cardinality see SCARD
func (storage *Client) Cardinality(key string) (int, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	return redis.Int(connection.Do("SCARD", key))
}

// AddToSet see SADD
func (storage *Client) AddToSet(key string, values ...[]byte) error {
	if len(values) == 0 {
		return nil
	}

	connection := storage.checkout()
	defer storage.release(connection)

	args := make([]interface{}, len(values)+1)
	args[0] = key
	for index, value := range values {
		args[index+1] = value
	}

	_, err := connection.Do("SADD", args...)
	return err
}

// RemoveFromSet see SREM
func (storage *Client) RemoveFromSet(key string, values ...[]byte) error {
	if len(values) == 0 {
		return nil
	}

	connection := storage.checkout()
	defer storage.release(connection)

	args := make([]interface{}, len(values)+1)
	args[0] = key
	for index, value := range values {
		args[index+1] = value
	}

	_, err := connection.Do("SREM", args...)
	return err
}

// GetAllFromSet see SMEMBERS
func (storage *Client) GetAllFromSet(key string) ([][]byte, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	data, err := redis.ByteSlices(connection.Do("SMEMBERS", key))
	if err == redis.ErrNil {
		return [][]byte{}, nil
	}
	return data, err
}

// IsMemberOfSet see SISMEMBER
func (storage *Client) IsMemberOfSet(key string, value []byte) (bool, error) {
	connection := storage.checkout()
	defer storage.release(connection)

	data, err := redis.Bool(connection.Do("SISMEMBER", key, value))
	return data, err
}

// StoreUnionSet see SUNIONSTORE
func (storage *Client) StoreUnionSet(key string, keys ...string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	connection := storage.checkout()
	defer storage.release(connection)

	args := make([]interface{}, len(keys)+1)
	args[0] = key
	for index, key := range keys {
		args[index+1] = key
	}

	return redis.Int(connection.Do("SUNIONSTORE", args...))
}

// Delete see DEL
func (storage *Client) Delete(keys ...string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	connection := storage.checkout()
	defer storage.release(connection)

	params := make([]interface{}, len(keys))
	for index, key := range keys {
		params[index] = key
	}

	count, err := redis.Int(connection.Do("DEL", params...))
	return count, err
}
