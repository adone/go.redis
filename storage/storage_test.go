package storage_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"

	"../pool"
	"../storage"
)

var _ = Describe("Client", func() {
	var (
		config storage.Configuration
		client *storage.Client

		connection *redigomock.Conn

		key   string = "foo"
		value []byte = []byte("bar")
	)

	BeforeEach(func() {
		connection = redigomock.NewConn()
	})

	DeleteTests := func() {
		var command *redigomock.Cmd

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("DEL", key).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				count, err := client.Delete(key)

				Expect(err).To(HaveOccurred())
				Expect(count).To(Equal(0))
			})
		})

		Context("succeed", func() {
			BeforeEach(func() {
				command = connection.Command("DEL", key).Expect(int64(1))
			})

			It("should delete value", func() {
				count, err := client.Delete(key)

				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})
	}

	IncrementTests := func() {
		var command *redigomock.Cmd

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("INCRBY", key, 1).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				count, err := client.Increment(key, 1)

				Expect(err).To(HaveOccurred())
				Expect(count).To(Equal(0))
			})
		})

		Context("succeed", func() {
			BeforeEach(func() {
				command = connection.Command("INCRBY", key, 1).Expect(int64(1))
			})

			It("should increment value", func() {
				count, err := client.Increment(key, 1)

				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})
	}

	ExpireTests := func() {
		var command *redigomock.Cmd

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("EXPIRE", key, 10).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				Expect(client.Expire(key, 10)).To(HaveOccurred())
			})
		})

		Context("success", func() {
			BeforeEach(func() {
				command = connection.Command("EXPIRE", key, 10).Expect(1)
			})

			It("should set key ttl", func() {
				Expect(client.Expire(key, 10*time.Second)).ToNot(HaveOccurred())
			})
		})
	}

	GetTests := func() {
		var command *redigomock.Cmd

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("GET", key).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				data, err := client.Get(key)

				Expect(err).To(HaveOccurred())
				Expect(data).To(BeEmpty())
			})
		})

		Context("when key exists", func() {
			BeforeEach(func() {
				command = connection.Command("GET", key).Expect(value)
			})

			It("should get value", func() {
				data, err := client.Get(key)

				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(Equal(value))
			})
		})

		Context("when key does not exist", func() {
			BeforeEach(func() {
				command = connection.Command("GET", key).ExpectError(redis.ErrNil)
			})

			It("should return empty value", func() {
				data, err := client.Get(key)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(BeEmpty())
			})
		})
	}

	MultiGetTests := func() {
		var (
			command *redigomock.Cmd
			key1    = "1"
			key2    = "2"
			key3    = "3"
			value1  = []byte("wut")
			value2  = []byte("a")
			value3  = []byte("fuck")
		)

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("MGET", key1, key2, key3).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				data, err := client.MultiGet(key1, key2, key3)

				Expect(err).To(HaveOccurred())
				Expect(data).To(BeEmpty())
			})
		})

		Context("when key exists", func() {
			BeforeEach(func() {
				command = connection.Command("MGET", key1, key2, key3).Expect([]interface{}{value1, value2, value3})
			})

			It("should get values", func() {
				data, err := client.MultiGet(key1, key2, key3)

				Expect(err).ToNot(HaveOccurred())
				Expect(data[0]).To(Equal(value1))
				Expect(data[1]).To(Equal(value2))
				Expect(data[2]).To(Equal(value3))
			})
		})

		Context("when one key does not exist", func() {
			BeforeEach(func() {
				command = connection.Command("MGET", key1, key2, key3).Expect([]interface{}{value1, nil, value3})
			})

			It("should return one empty value", func() {
				data, err := client.MultiGet(key1, key2, key3)
				Expect(err).ToNot(HaveOccurred())
				Expect(data[0]).To(Equal(value1))
				Expect(data[1]).To(BeEmpty())
				Expect(data[2]).To(Equal(value3))
			})
		})
	}

	PuplishTests := func() {
		var command *redigomock.Cmd

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("PUBLISH", key, value).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				Expect(client.Publish(key, value)).ToNot(Succeed())
			})
		})

		Context("succeed", func() {
			BeforeEach(func() {
				command = connection.Command("PUBLISH", key, value).Expect("ok")
			})

			It("should get value", func() {
				Expect(client.Publish(key, value)).To(Succeed())
			})
		})
	}

	SetFieldTests := func() {
		var (
			command *redigomock.Cmd
			field   string = "baz"
		)

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("HSET", key, field, value).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				Expect(client.SetField(key, field, value)).ToNot(Succeed())
			})
		})

		Context("succeed", func() {
			BeforeEach(func() {
				command = connection.Command("HSET", key, field, value).Expect("ok")
			})

			It("should set hash field", func() {
				Expect(client.SetField(key, field, value)).To(Succeed())
			})
		})
	}

	SetFieldsTests := func() {
		var (
			command1 *redigomock.Cmd
			command2 *redigomock.Cmd
			field1   = "one"
			field2   = "two"
			value1   = "A"
			value2   = 2
			hash     = map[string]interface{}{
				field1: value1,
				field2: value2,
			}
		)

		AfterEach(func() {
			Expect(1).To(SatisfyAny(Equal(connection.Stats(command1)), Equal(connection.Stats(command2))))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command1 = connection.Command("HMSET", key, field1, value1, field2, value2).
					ExpectError(fmt.Errorf("error"))
				command2 = connection.Command("HMSET", key, field2, value2, field1, value1).
					ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				Expect(client.SetFields(key, hash)).To(HaveOccurred())
			})
		})

		Context("succeed", func() {
			BeforeEach(func() {
				command1 = connection.Command("HMSET", key, field1, value1, field2, value2).
					Expect("ok")
				command2 = connection.Command("HMSET", key, field2, value2, field1, value1).
					Expect("ok")
			})

			It("should set hash field", func() {
				Expect(client.SetFields(key, hash)).ToNot(HaveOccurred())
			})
		})
	}

	GetFieldTests := func() {
		var (
			command *redigomock.Cmd
			field   string = "baz"
		)

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("HGET", key, field).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				data, err := client.GetField(key, field)
				Expect(err).To(HaveOccurred())
				Expect(data).To(BeEmpty())
			})
		})

		Context("when field exists", func() {
			BeforeEach(func() {
				command = connection.Command("HGET", key, field).Expect(value)
			})

			It("should return value", func() {
				data, err := client.GetField(key, field)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(Equal(value))
			})
		})

		Context("when field does not exist", func() {
			BeforeEach(func() {
				command = connection.Command("HGET", key, field).ExpectError(redis.ErrNil)
			})

			It("should return empty value", func() {
				data, err := client.GetField(key, field)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(BeEmpty())
			})
		})
	}

	HashIncrementTests := func() {
		var (
			command *redigomock.Cmd
			field   string = "baz"
		)

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("HINCRBY", key, field, 1).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				count, err := client.IncrementField(key, field, 1)

				Expect(err).To(HaveOccurred())
				Expect(count).To(Equal(0))
			})
		})

		Context("succeed", func() {
			BeforeEach(func() {
				command = connection.Command("HINCRBY", key, field, 1).Expect(int64(1))
			})

			It("should increment value", func() {
				count, err := client.IncrementField(key, field, 1)

				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})
	}

	FieldExistTests := func() {
		var (
			command *redigomock.Cmd
			field   string = "baz"
		)

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("HEXISTS", key, field).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				exist, err := client.FieldExist(key, field)
				Expect(err).To(HaveOccurred())
				Expect(exist).To(BeFalse())
			})
		})

		Context("when field exists", func() {
			BeforeEach(func() {
				command = connection.Command("HEXISTS", key, field).Expect(int64(1))
			})

			It("should return true", func() {
				exist, err := client.FieldExist(key, field)
				Expect(err).ToNot(HaveOccurred())
				Expect(exist).To(BeTrue())
			})
		})

		Context("when field does not exist", func() {
			BeforeEach(func() {
				command = connection.Command("HEXISTS", key, field).Expect(int64(0))
			})

			It("should return false", func() {
				exist, err := client.FieldExist(key, field)
				Expect(err).ToNot(HaveOccurred())
				Expect(exist).To(BeFalse())
			})
		})
	}

	GetValuesTests := func() {
		var (
			command *redigomock.Cmd
			value   []byte = []byte("baz")
		)

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("HVALS", key).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				data, err := client.GetValues(key)
				Expect(err).To(HaveOccurred())
				Expect(data).To(BeEmpty())
			})
		})

		Context("when key exists", func() {
			BeforeEach(func() {
				command = connection.Command("HVALS", key).Expect([]interface{}{value})
			})

			It("should return true", func() {
				data, err := client.GetValues(key)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(Equal([][]byte{value}))
			})
		})

		Context("when key does not exist", func() {
			BeforeEach(func() {
				command = connection.Command("HVALS", key).ExpectError(redis.ErrNil)
			})

			It("should return false", func() {
				data, err := client.GetValues(key)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(BeEmpty())
			})
		})
	}

	GetAllFromSetTests := func() {
		var (
			command *redigomock.Cmd
			value   []byte = []byte("baz")
		)

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("SMEMBERS", key).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				data, err := client.GetAllFromSet(key)
				Expect(err).To(HaveOccurred())
				Expect(data).To(BeEmpty())
			})
		})

		Context("when set exists", func() {
			BeforeEach(func() {
				command = connection.Command("SMEMBERS", key).Expect([]interface{}{value})
			})

			It("should return true", func() {
				data, err := client.GetAllFromSet(key)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(Equal([][]byte{value}))
			})
		})

		Context("when set does not exist", func() {
			BeforeEach(func() {
				command = connection.Command("SMEMBERS", key).ExpectError(redis.ErrNil)
			})

			It("should return false", func() {
				data, err := client.GetAllFromSet(key)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(BeEmpty())
			})
		})
	}

	IsMemberOfSetTests := func() {
		var (
			command *redigomock.Cmd
			value    []byte = []byte("baz")
		)

		AfterEach(func() {
			Expect(connection.Stats(command)).To(Equal(1))
		})

		Context("failed", func() {
			BeforeEach(func() {
				command = connection.Command("SISMEMBER", key, value).ExpectError(fmt.Errorf("error"))
			})

			It("should return error", func() {
				exist, err := client.IsMemberOfSet(key, value)
				Expect(err).To(HaveOccurred())
				Expect(exist).To(BeFalse())
			})
		})

		Context("when field exists", func() {
			BeforeEach(func() {
				command = connection.Command("SISMEMBER", key, value).Expect(int64(1))
			})

			It("should return true", func() {
				exist, err := client.IsMemberOfSet(key, value)
				Expect(err).ToNot(HaveOccurred())
				Expect(exist).To(BeTrue())
			})
		})

		Context("when field does not exist", func() {
			BeforeEach(func() {
				command = connection.Command("SISMEMBER", key, value).Expect(int64(0))
			})

			It("should return false", func() {
				exist, err := client.IsMemberOfSet(key, value)
				Expect(err).ToNot(HaveOccurred())
				Expect(exist).To(BeFalse())
			})
		})
	}

	Context("without connection", func() {
		BeforeEach(func() {
			config = storage.Configuration{}
		})

		It("should fail on creating storage client", func() {
			Expect(func() { storage.New(config) }).To(Panic())
		})
	})

	Context("with single connection", func() {
		BeforeEach(func() {
			config = storage.Configuration{
				Connection: connection,
			}
		})

		JustBeforeEach(func() {
			client = storage.New(config)
		})

		Describe("method GET", GetTests)
		Describe("method MGET", MultiGetTests)
		Describe("method EXPIRE", ExpireTests)
		Describe("method DELETE", DeleteTests)
		Describe("method INCREMENT", IncrementTests)
		Describe("method PUBLISH", PuplishTests)
		Describe("method SET FIELD", SetFieldTests)
		Describe("method SET FIELDS", SetFieldsTests)
		Describe("method FIELD INCREMENT", HashIncrementTests)
		Describe("method GET FIELD", GetFieldTests)
		Describe("method FIELD EXIST", FieldExistTests)
		Describe("method GET VALUES", GetValuesTests)
		Describe("method SISMEMBER", IsMemberOfSetTests)
		Describe("method SMEMBERS", GetAllFromSetTests)
	})

	Context("with connection pool", func() {
		BeforeEach(func() {
			config = storage.Configuration{
				Pool: pool.New(pool.Configuration{},
					func() (redis.Conn, error) { return connection, nil },
					pool.Check(pool.Configuration{}),
				),
			}
		})

		JustBeforeEach(func() {
			client = storage.New(config)
		})

		Describe("method GET", GetTests)
		Describe("method MGET", MultiGetTests)
		Describe("method EXPIRE", ExpireTests)
		Describe("method DELETE", DeleteTests)
		Describe("method INCREMENT", IncrementTests)
		Describe("method PUBLISH", PuplishTests)
		Describe("method SET FIELD", SetFieldTests)
		Describe("method SET FIELDS", SetFieldsTests)
		Describe("method FIELD INCREMENT", HashIncrementTests)
		Describe("method GET FIELD", GetFieldTests)
		Describe("method FIELD EXIST", FieldExistTests)
		Describe("method GET VALUES", GetValuesTests)
		Describe("method SISMEMBER", IsMemberOfSetTests)
		Describe("method SMEMBERS", GetAllFromSetTests)
	})
})
