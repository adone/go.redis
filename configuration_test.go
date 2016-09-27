package redis_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"../redis"
)

var _ = Describe("Configuration", func() {
	var (
		config         redis.Configuration
		commonTimeout  = 10 * time.Second
		connectTimeout = 1 * time.Second
		readTimeout    = 2 * time.Second
		writeTimeout   = 3 * time.Second
	)

	BeforeEach(func() {
		config = redis.Configuration{
			Timeout: commonTimeout,
		}
	})

	Context("without connect timeout", func() {
		It("connect timeout should be equal to commot timeout", func() {
			Expect(config.GetConnectTimeout()).To(Equal(commonTimeout))
		})
	})

	Context("with connect timeout", func() {
		JustBeforeEach(func() {
			config.ConnectTimeout = connectTimeout
		})

		It("connect timeout should be equal to commot timeout", func() {
			Expect(config.GetConnectTimeout()).To(Equal(connectTimeout))
		})
	})

	Context("without read timeout", func() {
		It("read timeout should be equal to commot timeout", func() {
			Expect(config.GetReadTimeout()).To(Equal(commonTimeout))
		})
	})

	Context("with read timeout", func() {
		JustBeforeEach(func() {
			config.ReadTimeout = readTimeout
		})

		It("read timeout should be equal to commot timeout", func() {
			Expect(config.GetReadTimeout()).To(Equal(readTimeout))
		})
	})

	Context("without write timeout", func() {
		It("write timeout should be equal to commot timeout", func() {
			Expect(config.GetWriteTimeout()).To(Equal(commonTimeout))
		})
	})

	Context("with write timeout", func() {
		JustBeforeEach(func() {
			config.WriteTimeout = writeTimeout
		})

		It("write timeout should be equal to commot timeout", func() {
			Expect(config.GetWriteTimeout()).To(Equal(writeTimeout))
		})
	})
})
