package engine

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

func Send(subject, content, from, to string) error {
	smtpServer := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	user := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	fromAddress := from
	toAddress := []string{to}

	tmpl, err := template.ParseFiles("templates/templates.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := map[string]string{
		"Subject": subject,
		"Content": content,
		"From":    from,
	}

	var body bytes.Buffer

	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	emailBody := fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", subject, body.String())

	auth := smtp.PlainAuth("", user, password, smtpServer)

	err = smtp.SendMail(
		smtpServer+":"+port,
		auth,
		fromAddress,
		toAddress,
		[]byte(emailBody),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
