package pool_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"../pool"
)

var _ = Describe("Environment", func() {
	var (
		prefix   = "TEST"
		prefixed = func(env string) string { return fmt.Sprintf("%s_%s", prefix, env) }
	)

	Context("method IdleTimeout", func() {
		It("should return empty timeout", func() {
			Expect(pool.IdleTimeout(prefix)).To(Equal(time.Duration(0)))
		})

		Context("when REDIS_POOL_TIMEOUT is set", func() {
			var (
				timeout = "1s"
				result  = time.Second
			)

			BeforeEach(func() {
				os.Setenv("REDIS_POOL_TIMEOUT", timeout)
			})

			AfterEach(func() {
				os.Setenv("REDIS_POOL_TIMEOUT", "")
			})

			It("should return idle timeout", func() {
				Expect(pool.IdleTimeout(prefix)).To(Equal(result))
			})

			Context("and prefixed REDIS_POOL_TIMEOUT is set", func() {
				var (
					ptimeout = "2s"
					presult  = 2 * time.Second
				)

				BeforeEach(func() {
					os.Setenv(prefixed("REDIS_POOL_TIMEOUT"), ptimeout)
				})

				AfterEach(func() {
					os.Setenv(prefixed("REDIS_POOL_TIMEOUT"), "")
				})

				It("should return idle timeout", func() {
					Expect(pool.IdleTimeout(prefix)).To(Equal(presult))
				})
			})
		})
	})

	Context("method CheckFrequency", func() {
		It("should return empty timeout", func() {
			Expect(pool.CheckFrequency(prefix)).To(Equal(time.Duration(0)))
		})

		Context("when REDIS_POOL_CHECK_TIMEOUT is set", func() {
			var (
				timeout = "1s"
				result  = time.Second
			)

			BeforeEach(func() {
				os.Setenv("REDIS_POOL_CHECK_TIMEOUT", timeout)
			})

			AfterEach(func() {
				os.Setenv("REDIS_POOL_CHECK_TIMEOUT", "")
			})

			It("should return idle timeout", func() {
				Expect(pool.CheckFrequency(prefix)).To(Equal(result))
			})

			Context("and prefixed REDIS_POOL_CHECK_TIMEOUT is set", func() {
				var (
					ptimeout = "2s"
					presult  = 2 * time.Second
				)

				BeforeEach(func() {
					os.Setenv(prefixed("REDIS_POOL_CHECK_TIMEOUT"), ptimeout)
				})

				AfterEach(func() {
					os.Setenv(prefixed("REDIS_POOL_CHECK_TIMEOUT"), "")
				})

				It("should return idle timeout", func() {
					Expect(pool.CheckFrequency(prefix)).To(Equal(presult))
				})
			})
		})
	})

	Context("method MaxIdleCount", func() {
		It("should return default count", func() {
			Expect(pool.MaxIdleCount(prefix)).To(Equal(pool.DefaultRedisPoolSize))
		})

		Context("when REDIS_IDLE_POOL_SIZE is set", func() {
			var (
				size   = "10"
				result = 10
			)

			BeforeEach(func() {
				os.Setenv("REDIS_IDLE_POOL_SIZE", size)
			})

			AfterEach(func() {
				os.Setenv("REDIS_IDLE_POOL_SIZE", "")
			})

			It("should return max idle connection count", func() {
				Expect(pool.MaxIdleCount(prefix)).To(Equal(result))
			})

			Context("and prefixed REDIS_IDLE_POOL_SIZE is set", func() {
				var (
					prefixedSize    = "20"
					prexexiedResult = 20
				)

				BeforeEach(func() {
					os.Setenv(prefixed("REDIS_IDLE_POOL_SIZE"), prefixedSize)
				})

				AfterEach(func() {
					os.Setenv(prefixed("REDIS_IDLE_POOL_SIZE"), "")
				})

				It("should return max idle connection count for concrete connection", func() {
					Expect(pool.MaxIdleCount(prefix)).To(Equal(prexexiedResult))
				})
			})
		})
	})

	Context("method MaxActiveCount", func() {
		It("should return default count", func() {
			Expect(pool.MaxActiveCount(prefix)).To(Equal(0))
		})

		Context("when REDIS_ACTIVE_POOL_SIZE is set", func() {
			var (
				size   = "10"
				result = 10
			)

			BeforeEach(func() {
				os.Setenv("REDIS_ACTIVE_POOL_SIZE", size)
			})

			AfterEach(func() {
				os.Setenv("REDIS_ACTIVE_POOL_SIZE", "")
			})

			It("should return max idle connection count", func() {
				Expect(pool.MaxActiveCount(prefix)).To(Equal(result))
			})

			Context("and prefixed REDIS_ACTIVE_POOL_SIZE is set", func() {
				var (
					prefixedSize    = "20"
					prexexiedResult = 20
				)

				BeforeEach(func() {
					os.Setenv(prefixed("REDIS_ACTIVE_POOL_SIZE"), prefixedSize)
				})

				AfterEach(func() {
					os.Setenv(prefixed("REDIS_ACTIVE_POOL_SIZE"), "")
				})

				It("should return max idle connection count for concrete connection", func() {
					Expect(pool.MaxActiveCount(prefix)).To(Equal(prexexiedResult))
				})
			})
		})
	})
})
