package user

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Describe("getEnv", func() {

		const (
			testEnvVarKey   = "test"
			testEnvVarValue = "test-value"
			defaultValue    = "default"
		)

		BeforeEach(func() {
			os.Setenv(testEnvVarKey, testEnvVarValue)
		})

		AfterEach(func() {
			os.Unsetenv(testEnvVarKey)
		})

		It("should return the env var value", func() {
			envVarValue := getEnv(testEnvVarKey, "")
			Expect(envVarValue).To(Equal(testEnvVarValue))
		})

		It("should return the default value if env var is not set", func() {
			os.Unsetenv(testEnvVarKey)

			envVarValue := getEnv(testEnvVarKey, defaultValue)
			Expect(envVarValue).To(Equal(defaultValue))
		})
	})

	Describe("NewRedisConfig", func() {

		BeforeEach(func() {
			os.Setenv(REDIS_ADDRESS, "")
			os.Setenv(REDIS_PASSWORD, "")
			os.Setenv(REDIS_DATABASE, "")
		})

		AfterEach(func() {
			os.Unsetenv(REDIS_ADDRESS)
			os.Unsetenv(REDIS_PASSWORD)
			os.Unsetenv(REDIS_DATABASE)
		})

		It("should return an error if REDIS_ADDRESS env var is missing", func() {
			os.Unsetenv(REDIS_ADDRESS)

			_, err := NewRedisConfig()
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if REDIS_PASSWORD env var is missing", func() {
			os.Unsetenv(REDIS_PASSWORD)

			_, err := NewRedisConfig()
			Expect(err).To(HaveOccurred())
		})

		It("should return no error and default redis instance if REDIS_DATABASE env var is missing", func() {
			os.Unsetenv(REDIS_DATABASE)

			config, err := NewRedisConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(config.GetDbInstance()).To(Equal(defaultRedisInstance))
		})

		It("should return no error if no env var is missing", func() {
			config, err := NewRedisConfig()

			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(Equal(redisConfig{}))
		})
	})

})
