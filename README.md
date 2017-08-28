# Redis

Helpers for [redis](github.com/garyburd/redigo/redis) package

## Usage

### Single connection

```go
  import "gopkg.in/adone/go.redis.v1"
```

```go
  config := redis.Configuration{
    Address:        "localhost:6379",
    ConnectTimeout: 1 * time.Second,
    ReadTimeout:    1 * time.Second,
    WriteTimeout:   5 * time.Second,
    Database:       1,
  }
```

```go
  // TEST_REDIS_ADDRESS=localhost:6379
  // TEST_REDIS_DATABASE=1
  // TEST_REDIS_TIMEOUT=1s
  // TEST_REDIS_WRITE_TIMEOUT=5s
  config := redis.ENV("TEST")
```

```go
  // TEST_REDIS_ADDRESS=localhost:6379
  // REDIS_TIMEOUT=1s
  // TEST_REDIS_TIMEOUT=2s

  config := redis.ENV("TEST")
  config.Timeout // => 2 * time.Second
```

```go
  // TEST_REDIS_SERVICE_HOST=10.0.0.2
  // TEST_REDIS_SERVICE_PORT=9736
  // TEST_REDIS_DATABASE=1

  config := redis.ENV("TEST")
  config.Address() // => "10.0.0.2:9736"
```


```go
  // TEST_REDIS_DATABASE=1
  redis.Database("TEST") // => 1

  // TEST_REDIS_TIMEOUT=1s
  redis.CommonTimeout("TEST") // => time.Second
```


```go
package main

import (
    "fmt"

    "gopkg.in/adone/go.redis.v1"
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

Support ENV variables:

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


##### PREFIX_SENTINEL_ADDRESSES
`PREFIX_SENTINEL_ADDRESSES=host1:port1,host2:port2,host3:port3`


### Pool

```go
  import "gopkg.in/adone/go.redis.v1/pool"
```

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
  config := pool.ENV("TEST")
```

Full example:

```go
package main

import (
    "fmt"

    redigo "github.com/garyburd/redigo/redis"

    "gopkg.in/adone/go.redis.v1"
    "gopkg.in/adone/go.redis.v1/pool"
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

Support ENV variables:

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

### Storage

```go
import "gopkg.in/adone/go.redis.v1/storage"

```

Commands:

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

`storage` automatic change SET to SETEX if `KeyTTL` is set

```go
  client.New(storage.Configuration{KeyTTL: 100})
  client.Set("key", []byte("value")) // => SETEX key 100 value
```

```go
  client.New(storage.Configuration{KeyTTL: func(key string) int { return 10 }})
  client.Set("key", []byte("value")) // => SETEX key 10 value
```

TTL can be set in `storage.Setter`

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

* SMEMBERS

```go
data, err := client.GetAllFromSet("setname")
```

* SISMEMBER

```go
exist, err := client.IsMemberOfSet("setname", "value")
```

```go
  client.DeleteFields("key", "field1", "filed2", "filed3")
```

* PUBLISH

```go
  err := client.Publish("key", []byte("value"))
```

Full example:

```go
package main

import (
    "fmt"

    "gopkg.in/adone/go.redis.v1"
    "gopkg.in/adone/go.redis.v1/pool"
    "gopkg.in/adone/go.redis.v1/storage"
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
