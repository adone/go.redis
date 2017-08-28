package redis_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"../redis"
)

var _ = Describe("Environment", func() {
	var (
		prefix   = "TEST"
		prefixed = func(env string) string { return fmt.Sprintf("%s_%s", prefix, env) }
	)

	Context("method Address", func() {
		AfterEach(func() {
			os.Setenv(prefixed("REDIS_SERVICE_HOST"), "")
			os.Setenv(prefixed("REDIS_SERVICE_PORT"), "")
			os.Setenv(prefixed("REDIS_PASSWORD"), "")
			os.Setenv(prefixed("REDIS_DATABASE"), "")
			os.Setenv(prefixed("REDIS_ADDRESS"), "")
		})

		Context("when SERVICE HOST is set", func() {
			var (
				host    = "example"
				address = "example:6379"
			)

			BeforeEach(func() {
				os.Setenv(prefixed("REDIS_SERVICE_HOST"), host)
			})

			It("should return address", func() {
				result, err := redis.ENV(prefix).Address()
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(address))
			})
		})

		Context("when SERVICE PORT is set", func() {
			var (
				port    = "9736"
				address = "localhost:9736"
			)

			BeforeEach(func() {
				os.Setenv(prefixed("REDIS_SERVICE_PORT"), port)
			})

			It("should return address", func() {
				result, err := redis.ENV(prefix).Address()
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(address))
			})
		})

		Context("when SERVICE HOST & PORT are set", func() {
			var (
				port    = "9736"
				host    = "lvh.me"
				address = "lvh.me:9736"
			)

			BeforeEach(func() {
				os.Setenv(prefixed("REDIS_SERVICE_PORT"), port)
				os.Setenv(prefixed("REDIS_SERVICE_HOST"), host)
			})

			It("should return address", func() {
				result, err := redis.ENV(prefix).Address()
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(address))
			})
		})

		Context("when REDIS_ADDRESS is set", func() {
			var address = "127.0.0.1:6379"

			BeforeEach(func() {
				os.Setenv(prefixed("REDIS_ADDRESS"), address)
			})

			It("should return address", func() {
				result, err := redis.ENV(prefix).Address()
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(address))
			})
		})
	})

	Context("method CommonTimeout", func() {
		It("should return empty timeout", func() {
			Expect(redis.CommonTimeout(prefix)).To(Equal(time.Duration(0)))
		})

		Context("when REDIS_TIMEOUT is set", func() {
			var (
				timeout = "1s"
				result  = time.Second
			)

			BeforeEach(func() {
				os.Setenv("REDIS_TIMEOUT", timeout)
			})

			AfterEach(func() {
				os.Setenv("REDIS_TIMEOUT", "")
			})

			It("should return redis timeout", func() {
				Expect(redis.CommonTimeout(prefix)).To(Equal(result))
			})

			Context("and prefixed REDIS_TIMEOUT is set", func() {
				var (
					ptimeout = "2s"
					presult  = 2 * time.Second
				)

				BeforeEach(func() {
					os.Setenv(prefixed("REDIS_TIMEOUT"), ptimeout)
				})

				AfterEach(func() {
					os.Setenv(prefixed("REDIS_TIMEOUT"), "")
				})

				It("should return redis timeout", func() {
					Expect(redis.CommonTimeout(prefix)).To(Equal(presult))
				})
			})
		})
	})

	Context("method ConnectionTimeout", func() {
		It("should return empty timeout", func() {
			Expect(redis.ConnectionTimeout(prefix)).To(Equal(time.Duration(0)))
		})

		Context("when REDIS_CONNECT_TIMEOUT is set", func() {
			var (
				timeout = "1s"
				result  = time.Second
			)

			BeforeEach(func() {
				os.Setenv("REDIS_CONNECT_TIMEOUT", timeout)
			})

			AfterEach(func() {
				os.Setenv("REDIS_CONNECT_TIMEOUT", "")
			})

			It("should return redis timeout", func() {
				Expect(redis.ConnectionTimeout(prefix)).To(Equal(result))
			})

			Context("and prefixed REDIS_CONNECT_TIMEOUT is set", func() {
				var (
					ptimeout = "2s"
					presult  = 2 * time.Second
				)

				BeforeEach(func() {
					os.Setenv(prefixed("REDIS_CONNECT_TIMEOUT"), ptimeout)
				})

				AfterEach(func() {
					os.Setenv(prefixed("REDIS_CONNECT_TIMEOUT"), "")
				})

				It("should return redis timeout", func() {
					Expect(redis.ConnectionTimeout(prefix)).To(Equal(presult))
				})
			})
		})
	})

	Context("method ReadTimeout", func() {
		It("should return empty timeout", func() {
			Expect(redis.ReadTimeout(prefix)).To(Equal(time.Duration(0)))
		})

		Context("when REDIS_READ_TIMEOUT is set", func() {
			var (
				timeout = "1s"
				result  = time.Second
			)

			BeforeEach(func() {
				os.Setenv("REDIS_READ_TIMEOUT", timeout)
			})

			AfterEach(func() {
				os.Setenv("REDIS_READ_TIMEOUT", "")
			})

			It("should return redis timeout", func() {
				Expect(redis.ReadTimeout(prefix)).To(Equal(result))
			})

			Context("and prefixed REDIS_READ_TIMEOUT is set", func() {
				var (
					ptimeout = "2s"
					presult  = 2 * time.Second
				)

				BeforeEach(func() {
					os.Setenv(prefixed("REDIS_READ_TIMEOUT"), ptimeout)
				})

				AfterEach(func() {
					os.Setenv(prefixed("REDIS_READ_TIMEOUT"), "")
				})

				It("should return redis timeout", func() {
					Expect(redis.ReadTimeout(prefix)).To(Equal(presult))
				})
			})
		})
	})

	Context("method WriteTimeout", func() {
		It("should return empty timeout", func() {
			Expect(redis.WriteTimeout(prefix)).To(Equal(time.Duration(0)))
		})

		Context("when REDIS_WRITE_TIMEOUT is set", func() {
			var (
				timeout = "1s"
				result  = time.Second
			)

			BeforeEach(func() {
				os.Setenv("REDIS_WRITE_TIMEOUT", timeout)
			})

			AfterEach(func() {
				os.Setenv("REDIS_WRITE_TIMEOUT", "")
			})

			It("should return redis timeout", func() {
				Expect(redis.WriteTimeout(prefix)).To(Equal(result))
			})

			Context("and prefixed REDIS_WRITE_TIMEOUT is set", func() {
				var (
					ptimeout = "2s"
					presult  = 2 * time.Second
				)

				BeforeEach(func() {
					os.Setenv(prefixed("REDIS_WRITE_TIMEOUT"), ptimeout)
				})

				AfterEach(func() {
					os.Setenv(prefixed("REDIS_WRITE_TIMEOUT"), "")
				})

				It("should return redis timeout", func() {
					Expect(redis.WriteTimeout(prefix)).To(Equal(presult))
				})
			})
		})
	})
})
