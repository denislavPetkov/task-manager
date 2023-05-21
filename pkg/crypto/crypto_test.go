package crypto

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crypto", func() {

	Describe("GenerateToken", func() {
		It("should return a random token", func() {
			token, err := GenerateToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())
		})
	})

	Describe("GetHashedPassword", func() {
		It("should return the password hashed", func() {
			hashedPassword, err := GetHashedPassword("test")
			Expect(err).NotTo(HaveOccurred())
			Expect(hashedPassword).NotTo(BeEmpty())
		})
	})

	Describe("IsHashedPasswordCorrect", func() {
		var (
			password       = "test"
			hashedPassword string
			err            error
		)

		BeforeEach(func() {
			hashedPassword, err = GetHashedPassword(password)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil when the passwords are the same", func() {
			err := IsHashedPasswordCorrect(password, hashedPassword)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error when the passwords are not the same", func() {
			err := IsHashedPasswordCorrect("wrong password", hashedPassword)
			Expect(err).To(HaveOccurred())
		})
	})
})
