package email

import (
	"errors"
	"net/smtp"

	smtpfacade "github.com/denislavpetkov/task-manager/pkg/facade/net/smtp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Email", func() {

	var (
		mockSmtpInstance *mockSmtp
	)

	BeforeEach(func() {
		mockSmtpInstance = &mockSmtp{}

		smtpfacade.SetSmtpInstance(mockSmtpInstance)
	})

	Describe("SendRecoveryEmail", func() {
		It("should return an error when an error occurs when sending email", func() {
			mockSmtpInstance.err = errors.New("")

			err := SendRecoveryEmail("", "")
			Expect(err).To(HaveOccurred())
		})

		It("should return nil when an error doesn not occur when sending email", func() {
			mockSmtpInstance.err = nil

			err := SendRecoveryEmail("", "")
			Expect(err).NotTo(HaveOccurred())
		})
	})

})

type mockSmtp struct {
	err error
}

func (m *mockSmtp) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return m.err
}
