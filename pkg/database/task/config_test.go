package task

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Describe("NewMongodbConfig", func() {

		BeforeEach(func() {
			os.Setenv(MONGODB_CONNECTION_URI, "test")
			os.Setenv(MONGODB_DATABASE_NAME, "test")
		})

		AfterEach(func() {
			os.Unsetenv(MONGODB_CONNECTION_URI)
			os.Unsetenv(MONGODB_DATABASE_NAME)
		})

		It("should return an error if MONGODB_CONNECTION_URI env var is missing", func() {
			os.Unsetenv(MONGODB_CONNECTION_URI)

			_, err := NewMongodbConfig()
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if MONGODB_DATABASE_NAME env var is missing", func() {
			os.Unsetenv(MONGODB_DATABASE_NAME)

			_, err := NewMongodbConfig()
			Expect(err).To(HaveOccurred())
		})

		It("should return no error if no env var is missing", func() {
			config, err := NewMongodbConfig()

			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(Equal(mongodbConfig{}))
		})
	})

})
