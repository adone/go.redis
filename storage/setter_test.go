package storage_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rafaeljusto/redigomock"

	"../storage"
)

var _ = Describe("Setter", func() {
	var (
		key   = "foo"
		value = []byte("bar")

		ttl        interface{}
		setter     storage.Setter
		command    *redigomock.Cmd
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
		setter = storage.Setter{
			TTL:     ttl,
			Storage: client,
			Key:     key,
			Value:   value,
		}
	})

	AfterEach(func() {
		Expect(connection.Stats(command)).To(Equal(1))
	})

	Context("when TTL is undefined", func() {
		BeforeEach(func() {
			ttl = nil
			command = connection.Command("SET", key, value).Expect("ok")
		})

		It("should set value without expire option", func() {
			Expect(setter.Call()).To(Succeed())
		})
	})

	Context("when TTL is function", func() {
		var expire = 10

		BeforeEach(func() {
			ttl = func(string) int { return expire }
			command = connection.Command("SETEX", key, expire, value).Expect("ok")
		})

		It("should set value with expire option", func() {
			Expect(setter.Call()).To(Succeed())
		})
	})

	Context("when TTL is int", func() {
		var expire = 20

		BeforeEach(func() {
			ttl = expire
			command = connection.Command("SETEX", key, expire, value).Expect("ok")
		})

		It("should set value with expire option", func() {
			Expect(setter.Call()).To(Succeed())
		})
	})

	Context("when TTL is duration", func() {
		var expire = 1

		BeforeEach(func() {
			ttl = time.Second
			command = connection.Command("SETEX", key, expire, value).Expect("ok")
		})

		It("should set value with expire option", func() {
			Expect(setter.Call()).To(Succeed())
		})
	})

	Context("manual set", func() {
		var expire = 1

		It("should set value without expire option", func() {
			command = connection.Command("SET", key, value).Expect("ok")
			Expect(setter.Set(0)).To(Succeed())
		})

		It("should set value without expire option", func() {
			command = connection.Command("SETEX", key, expire, value).Expect("ok")
			Expect(setter.Set(expire)).To(Succeed())
		})

		It("should return error", func() {
			command = connection.Command("SETEX", key, expire, value).ExpectError(fmt.Errorf("error"))
			Expect(setter.Set(expire)).ToNot(Succeed())
		})
	})
})
