# Redis

Данная библиотека упрощает работу с [redis](github.com/garyburd/redigo/redis):

* добавлены конфигурации для redis.Conn и redis.Pool
* возможность конфигурации через ENV
* добавлен тип Storage для более удобной работы с данными в Redis
 * обход ключей по шаблону c помощью SCAN
 * гибкая настройка TTL ключей
 * полная поддержка операций с типом Hash

## Использование в приложениях

### Обычное подключение

```go
  import "git.life-team.net/libs/redis"
```

Создается стандартный `redis.Conn`.

Конфигурация подключения возможна как вручную, так и черер ENV параметры:

```go
  config := redis.Configuration{
    URL:            "redis://localhost:6379",
    ConnectTimeout: 1 * time.Second,
    ReadTimeout:    1 * time.Second,
    WriteTimeout:   5 * time.Second,
    Database:       1,
  }
```

```go
  // TEST_REDIS_URL=redis://localhost:6379/1
  // TEST_REDIS_TIMEOUT=1s
  // TEST_REDIS_WRITE_TIMEOUT=5s
  config := redis.ENV("TEST")
```

При конфигурации через ENV параметры с префиксом перекрываются общие:

```go
  // TEST_REDIS_URL=redis://localhost:6379
  // REDIS_TIMEOUT=1s
  // TEST_REDIS_TIMEOUT=2s

  config := redis.ENV("TEST")
  config.Timeout // => 2 * time.Second
```

URL может быть собран через отдельные элементы, например, при работе внутри kubernates:

```go
  // TEST_REDIS_SERVICE_HOST=10.0.0.2
  // TEST_REDIS_SERVICE_PORT=9736
  // TEST_REDIS_DATABASE=1

  config := redis.ENV("TEST")
  config.URL // => "redis://10.0.0.2:9736/1"
```

Также возможно извлекать из ENV отдельные аргументы:

```go
  // TEST_REDIS_URL=redis://localhost:6379/1
  redis.URL("TEST") // => "redis://localhost:6379/1"

  // TEST_REDIS_TIMEOUT=1s
  redis.CommonTimeout("TEST") // => time.Second
```

Полный пример создания и использования подключения:

```go
package main

import (
    "fmt"

    "git.life-team.net/libs/redis"
)

func main() {
    conn, err := redis.New(redis.ENV("TEST"))
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    message, err := redis.String(conn.Do("PING", "FOOBAR"))
    fmt.Println(message, err)
}
```

Поддерживаемые ENV параметры:

* PREFIX_REDIS_URL
* PREFIX_REDIS_SERVICE_HOST
* PREFIX_REDIS_SERVICE_PORT
* PREFIX_REDIS_USERNAME
* PREFIX_REDIS_PASSWORD
* PREFIX_REDIX_DATABASE
* PREFIX_REDIS_TIMEOUT
* REDIS_TIMEOUT
* PREFIX_REDIS_CONNECT_TIMEOUT
* REDIS_CONNECT_TIMEOUT
* PREFIX_REDIS_WRITE_TIMEOUT
* REDIS_WRITE_TIMEOUT
* PREFIX_REDIS_READ_TIMEOUT
* REDIS_READ_TIMEOUT
* PREFIX_SENTINEL_HOSTS_PORTS 
* PREFIX_SENTINEL_MASTER_NAME


##### PREFIX_SENTINEL_HOSTS_PORTS
`PREFIX_SENTINEL_HOSTS_PORTS=sentinel-1:port1,sentinel-1:port2,sentinel-2:port1`


```

### Пулл подключений

```go
  import "git.life-team.net/libs/redis/pool"
```

Конфигурация аналогична обычному подключению - можно создать вручную, можно загрузить из ENV:

```go
  config := pool.Configuration{
    WaitConnection:           true, 
    MaxIdleConnectionCount:   8,
    MaxActiveConnectionCount: 32,
    IdleConnectionTimeout:    1 * time.Hour,
    CheckConnectionFrequency: 1 * time.Minute,
  }
```

```go
  // TEST_REDIS_ACTIVE_POOL_SIZE=32
  // TEST_REDIS_IDLE_POOL_SIZE=8
  // TEST_REDIS_POOL_TIMEOUT=1h
  // TEST_REDIS_POOL_CHECK_TIMEOUT=1m
  config := redis.ENV("TEST")
```

Полный пример использования:

```go
package main

import (
    "fmt"
    
    redigo "github.com/garyburd/redigo/redis"

    "git.life-team.net/libs/redis"
    "git.life-team.net/libs/redis/pool"
)

func main() {
    config := pool.ENV("TEST")
    
    pl := pool.New(config,
        redis.Connect(redis.ENV("TEST")),
        pool.Check(config),
    )

    conn := pl.Get()
    defer conn.Close()

    message, err := redis.String(conn.Do("PING", "FOOBAR"))
    fmt.Println(message, err)
}
```

Поддерживаемые ENV параметры:

* PREFIX_REDIS_ACTIVE_POOL_SIZE
* REDIS_ACTIVE_POOL_SIZE
* PREFIX_REDIS_IDLE_POOL_SIZE
* REDIS_IDLE_POOL_SIZE
* REDIS_POOL_SIZE
* PREFIX_REDIS_POOL_TIMEOUT
* REDIS_POOL_IDLE_TIMEOUT
* REDIS_POOL_TIMEOUT
* PREFIX_REDIS_POOL_CHECK_TIMEOUT
* REDIS_POOL_CHECK_TIMEOUT

### Хранилище

```go
import "git.life-team.net/libs/redis/storage"

```

Поддерживаемые команды:

* INCRBY

```go
  err := client.Increment("key", 1)
```

* EXPIRE

```go
  err := client.Expire("key", 1)
```

```go
  err := client.Expire("key", time.Minute)
```

```go
  err := client.Expire("key", func(key string) int { return 1 })
```

* SET

```go
  err := client.Set("key", []byte("value"))
```

* SETEX

`storage` автоматически меняет команду SET на SETEX, если установлена опция `KeyTTL`

```go
  client.New(storage.Configuration{KeyTTL: 100})
  client.Set("key", []byte("value")) // => SETEX key 100 value
```

```go
  client.New(storage.Configuration{KeyTTL: func(key string) int { return 10 }})
  client.Set("key", []byte("value")) // => SETEX key 10 value
```

TTL может настраиваться вручуню через тип `storage.Setter`

```go
  setter := storage.Setter{
    Storage: client,
    Key:     "key",
    Value:   []byte("value"),
    TTL:     1*time.Hour,
  }
  err := setter.Call()
```

```go
  setter := storage.Setter{
    Storage: client,
    Key:     "key",
    Value:   []byte("value"),
  }
  err := setter.Set(3600)
```

* GET

```go
  value, err := client.Get("key")
```

* DEL

```go
  err := client.Delete("key")
```

```go
  err := client.Delete("key1", "key2", "key3")
```

* SCAN

```go
  keys, err := client.Keys("key.*.template")
```

более гибкую настройку вызова обеспечивает тип `storage.Iterator`

```go
  iterator := storage.Iterator{
    Storage:   client,
    Template:  "key.*.template",
    BatchSize: 100,
  }

  keys, err := iterator.Call()
```

* HSET

```go
  client.AddField("key", "field", []byte("value"))
```

* HEXIST

```go
  client.FieldExist("key", "field")
```

* HGET

```go
  client.GetField("key", "value")
```

* HVALS

```go
  values, err := client.GetValues("key")
```

* HDEL

```go
  client.DeleteFields("key", "field")
```

поддерживается множественное удаление ключей в Hash

```go
  client.DeleteFields("key", "field1", "filed2", "filed3")
```

* PUBLISH

```go
  err := client.Publish("key", []byte("value"))
```

Полный пример использования:

```
package main

import (
    "fmt"

    "git.life-team.net/libs/redis"
    "git.life-team.net/libs/redis/pool"
    "git.life-team.net/libs/redis/storage"
)

func main() {
    config := storage.ENV("TEST")
    config.Pool = pool.New(pool.ENV("TEST"),
        redis.Connect(redis.ENV("TEST")),
        pool.Check(pool.ENV("TEST")),
    )

    client = storage.New(config)
    
    err := client.Set("foo", []byte("bar"))
    fmt.Println(err)
    
    message, err := client.Get("foo")
    fmt.Printf("%s %v", message, err)
    
    count, err := client.Delete("foo")
    fmt.Println(count, err)
}
```
