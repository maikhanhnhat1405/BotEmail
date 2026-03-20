package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

func SendReply(to, subject, body, host, port, user, pass string) error {
	auth := smtp.PlainAuth("", user, pass, host)

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\r\n", user))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return smtp.SendMail(host+":"+port, auth, user, []string{to}, []byte(msg.String()))
}