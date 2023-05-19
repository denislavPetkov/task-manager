package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

const (
	gmailHost     = "smtp.gmail.com"
	gmailHostPort = ":587"
)

var (
	serviceAccountUsername string
	serviceAccountPassword string
)

func init() {
	serviceAccountUsername = os.Getenv("EMAIL_USERNAME")
	serviceAccountPassword = os.Getenv("EMAIL_PASSWORD")
}

func SendRecoveryEmail(recipient, url string) error {

	msg := "From: " + serviceAccountUsername + "\n" +
		"To: " + recipient + "\n" +
		"Subject: Password Recovery\n\n" +
		fmt.Sprintf("Set a new password here: %s", url)

	err := smtp.SendMail(fmt.Sprintf("%s%s", gmailHost, gmailHostPort),
		smtp.PlainAuth("", serviceAccountUsername, serviceAccountPassword, gmailHost),
		serviceAccountUsername, []string{recipient}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	return nil
}
