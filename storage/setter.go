package storage

type Setter struct {
	Storage *Client
	TTL     interface{}
	Key     string
	Value   []byte
}

func (setter Setter) Call() error {
	return setter.Set(TTL{
		Key:   setter.Key,
		Value: setter.TTL,
	}.Seconds())
}

func (setter Setter) Set(ttl int) error {
	connection := setter.Storage.checkout()
	defer setter.Storage.release(connection)

	if ttl == 0 {
		_, err := connection.Do("SET", setter.Key, setter.Value)
		return err
	}

	_, err := connection.Do("SETEX", setter.Key, ttl, setter.Value)
	return err
}
