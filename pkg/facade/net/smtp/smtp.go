package smtp

import (
	originalsmtp "net/smtp"
)

type Smtp interface {
	SendMail(addr string, a originalsmtp.Auth, from string, to []string, msg []byte) error
}

type smtp struct{}

var smtpInstance Smtp = &smtp{}

func SetSmtpInstance(smtp Smtp) {
	smtpInstance = smtp
}

func GetSmtpInstance() Smtp {
	return smtpInstance
}

func (s *smtp) SendMail(addr string, a originalsmtp.Auth, from string, to []string, msg []byte) error {
	return originalsmtp.SendMail(addr, a, from, to, msg)
}
