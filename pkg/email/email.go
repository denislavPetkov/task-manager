package email

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	smtpfacade "github.com/denislavpetkov/task-manager/pkg/facade/net/smtp"
	"go.uber.org/zap"
)

const (
	gmailHost     = "smtp.gmail.com"
	gmailHostPort = ":587"
)

var (
	serviceAccountUsername string
	serviceAccountPassword string

	logger *zap.Logger
)

func init() {
	logger = zap.L().Named("email")
}

func init() {
	serviceAccountUsername = os.Getenv("EMAIL_USERNAME")
	serviceAccountPassword = os.Getenv("EMAIL_PASSWORD")
}

func SendRecoveryEmail(recipient, url string) error {

	msg := "From: " + serviceAccountUsername + "\n" +
		"To: " + recipient + "\n" +
		"Subject: Password Recovery\n\n" +
		fmt.Sprintf("Set a new password here: %s\nYou have %v before the link expires.", url, constants.PasswordRecoveryTokenExpiration)

	err := smtpfacade.GetSmtpInstance().SendMail(fmt.Sprintf("%s%s", gmailHost, gmailHostPort),
		smtp.PlainAuth("", serviceAccountUsername, serviceAccountPassword, gmailHost),
		serviceAccountUsername, []string{recipient}, []byte(msg))

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send an email, error: %v", err))
		return err
	}

	return nil
}
