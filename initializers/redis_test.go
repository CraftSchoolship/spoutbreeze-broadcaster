package initializers_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"spoutbreeze/initializers"
)

var _ = Describe("Redis Connection", func() {
	Describe("ConnectToRedis function", func() {
		Context("when all required environment variables are missing", func() {
			It("should panic with missing REDIS_PASSWORD", func() {
		
				originalPassword := os.Getenv("REDIS_PASSWORD")
				originalHost := os.Getenv("REDIS_HOST")
				originalPort := os.Getenv("REDIS_PORT")

				os.Unsetenv("REDIS_PASSWORD")
				os.Unsetenv("REDIS_HOST")
				os.Unsetenv("REDIS_PORT")

				defer func() {

					if originalPassword != "" {
						os.Setenv("REDIS_PASSWORD", originalPassword)
					}
					if originalHost != "" {
						os.Setenv("REDIS_HOST", originalHost)
					}
					if originalPort != "" {
						os.Setenv("REDIS_PORT", originalPort)
					}
				}()

				Expect(func() {
					initializers.ConnectToRedis()
				}).To(Panic())
			})
		})

		Context("when REDIS_PASSWORD is set but REDIS_HOST is missing", func() {
			It("should panic with missing REDIS_HOST", func() {
				// Store original values
				originalPassword := os.Getenv("REDIS_PASSWORD")
				originalHost := os.Getenv("REDIS_HOST")
				originalPort := os.Getenv("REDIS_PORT")


				os.Setenv("REDIS_PASSWORD", "test-password")
				os.Unsetenv("REDIS_HOST")
				os.Unsetenv("REDIS_PORT")

				defer func() {
					
					if originalPassword != "" {
						os.Setenv("REDIS_PASSWORD", originalPassword)
					} else {
						os.Unsetenv("REDIS_PASSWORD")
					}
					if originalHost != "" {
						os.Setenv("REDIS_HOST", originalHost)
					}
					if originalPort != "" {
						os.Setenv("REDIS_PORT", originalPort)
					}
				}()

				Expect(func() {
					initializers.ConnectToRedis()
				}).To(Panic())
			})
		})

		Context("when REDIS_PASSWORD and REDIS_HOST are set but REDIS_PORT is missing", func() {
			It("should panic with missing REDIS_PORT", func() {
				
				originalPassword := os.Getenv("REDIS_PASSWORD")
				originalHost := os.Getenv("REDIS_HOST")
				originalPort := os.Getenv("REDIS_PORT")

				
				os.Setenv("REDIS_PASSWORD", "test-password")
				os.Setenv("REDIS_HOST", "localhost")
				os.Unsetenv("REDIS_PORT")

				defer func() {
					
					if originalPassword != "" {
						os.Setenv("REDIS_PASSWORD", originalPassword)
					} else {
						os.Unsetenv("REDIS_PASSWORD")
					}
					if originalHost != "" {
						os.Setenv("REDIS_HOST", originalHost)
					} else {
						os.Unsetenv("REDIS_HOST")
					}
					if originalPort != "" {
						os.Setenv("REDIS_PORT", originalPort)
					}
				}()

				Expect(func() {
					initializers.ConnectToRedis()
				}).To(Panic())
			})
		})
	})

	Describe("Redis context and client", func() {
		It("should have a valid Redis context", func() {
			context := initializers.RedisContext
			Expect(context).NotTo(BeNil())
		})

		It("should have Redis client variable accessible", func() {		
			client := initializers.RedisClient
			_ = client // Just accessing it to cover the variable
		})
	})
})

func TestRedis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis Initializer Suite")
}
