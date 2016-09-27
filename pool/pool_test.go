package pool_test

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"../pool"
)

var _ = Describe("Connection", func() {
	var (
		config         pool.Configuration
		connection     *redigomock.Conn
		connectionPool *redis.Pool
	)

	BeforeEach(func() {
		connection = redigomock.NewConn()
	})

	JustBeforeEach(func() {
		connectionPool = pool.New(config,
			func() (redis.Conn, error) {
				return connection, nil
			},
			pool.Check(config),
		)
	})

	Context("when connection check event occurred", func() {
		var command *redigomock.Cmd

		BeforeEach(func() {
			config.MaxIdleConnectionCount = 1
			config.CheckConnectionFrequency = 10 * time.Millisecond
			command = connection.Command("PING").Expect("PONG")
		})

		AfterEach(func() {
			Expect(connection.Stats(command)).To(BeNumerically(">", 0))
		})

		It("should send ping", func() {
			connectionPool.Get().Close()
			time.Sleep(15 * time.Millisecond)
			connectionPool.Get().Close()
		})
	})

})
