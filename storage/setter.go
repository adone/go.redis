package storage

// Setter структура для записи данных в Redis
type Setter struct {
	Storage *Client     // соединение к редису
	TTL     interface{} // время жизни ключа
	Key     string      // название ключа
	Value   []byte      // данные
}

// Call запись данных в Redis
func (setter Setter) Call() error {
	return setter.Set(TTL{
		Key:   setter.Key,
		Value: setter.TTL,
	}.Seconds())
}

// Set запись данных в Redis на ttl секунд
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
