package storage_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"

	"../storage"
)

var _ = Describe("Iterator", func() {
	var (
		template string
		size     int

		iterator   *storage.Iterator
		commands   []*redigomock.Cmd
		connection *redigomock.Conn
		client     *storage.Client
	)

	BeforeEach(func() {
		connection = redigomock.NewConn()
		client = storage.New(storage.Configuration{
			Connection: connection,
		})
	})

	JustBeforeEach(func() {
		iterator = storage.NewIterator(
			storage.WithStorage(client), storage.WithBatchSize(size), storage.WithTemplate(template),
		)
	})

	Context("when BatchSize is undefined", func() {
		BeforeEach(func() {
			size = 0
		})

		It("should return empty array", func() {
			err := iterator.All(func(values []interface{}) {
				Expect(values).To(BeEmpty())
			})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when Template is undefined", func() {
		BeforeEach(func() {
			size = 2
			template = ""
			params := map[string][]interface{}{
				"1": {"SCAN", "0", "COUNT", size, []byte("foo1"), []byte("bar1")},
				"2": {"SCAN", "1", "COUNT", size, []byte("foo2"), []byte("bar2")},
				"0": {"SCAN", "2", "COUNT", size, []byte("foo3"), []byte("bar3")},
			}

			commands = make([]*redigomock.Cmd, 0, len(params))
			for cursor, args := range params {
				commands = append(commands,
					connection.Command(args[0].(string), args[1:4]...).Expect([]interface{}{[]byte(cursor), args[4:]}),
				)
			}
		})

		It("should return all keys", func() {
			var keys []string
			err := iterator.All(func(values []interface{}) {
				if found, err := redis.Strings(values, nil); err == nil {
					keys = append(keys, found...)
				}
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(keys).To(Equal([]string{"foo1", "bar1", "foo2", "bar2", "foo3", "bar3"}))
		})

		AfterEach(func() {
			for _, command := range commands {
				Expect(connection.Stats(command)).To(Equal(1))
			}
		})
	})

	Context("with Template", func() {
		BeforeEach(func() {
			size = 2
			template = "test"
			params := map[string][]interface{}{
				"1": {"SCAN", "0", "MATCH", template, "COUNT", size, []byte("foo1"), []byte("bar1")},
				"2": {"SCAN", "1", "MATCH", template, "COUNT", size, []byte("foo2"), []byte("bar2")},
				"0": {"SCAN", "2", "MATCH", template, "COUNT", size, []byte("foo3"), []byte("bar3")},
			}

			commands = make([]*redigomock.Cmd, 0, len(params))
			for cursor, args := range params {
				commands = append(commands,
					connection.Command(args[0].(string), args[1:6]...).Expect([]interface{}{[]byte(cursor), args[6:]}),
				)
			}
		})

		It("should return all keys", func() {
			var keys []string
			err := iterator.All(func(values []interface{}) {
				if found, err := redis.Strings(values, nil); err == nil {
					keys = append(keys, found...)
				}
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(keys).To(Equal([]string{"foo1", "bar1", "foo2", "bar2", "foo3", "bar3"}))
		})

		AfterEach(func() {
			for _, command := range commands {
				Expect(connection.Stats(command)).To(Equal(1))
			}
		})
	})

	Context("when error occurred", func() {
		BeforeEach(func() {
			size = 2
			template = "test"
			params := map[string][]interface{}{
				"1":     {"SCAN", "0", "MATCH", template, "COUNT", size, []byte("foo1"), []byte("bar1")},
				"2":     {"SCAN", "1", "MATCH", template, "COUNT", size, []byte("foo2"), []byte("bar2")},
				"error": {"SCAN", "2", "MATCH", template, "COUNT", size},
			}

			commands = make([]*redigomock.Cmd, 0, len(params))
			for cursor, args := range params {
				if cursor == "error" {
					connection.Command(args[0].(string), args[1:6]...).ExpectError(fmt.Errorf(cursor))
					continue
				}

				commands = append(commands,
					connection.Command(args[0].(string), args[1:6]...).Expect([]interface{}{[]byte(cursor), args[6:]}),
				)
			}
		})

		It("should return no keys", func() {
			var keys []string
			err := iterator.All(func(values []interface{}) {
				if found, err := redis.Strings(values, nil); err == nil {
					keys = append(keys, found...)
				}
			})
			Expect(err).To(HaveOccurred())
			Expect(keys).To(Equal([]string{"foo1", "bar1", "foo2", "bar2"}))
		})

		AfterEach(func() {
			for _, command := range commands {
				Expect(connection.Stats(command)).To(Equal(1))
			}
		})
	})
})
