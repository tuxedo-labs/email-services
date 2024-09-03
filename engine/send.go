package engine

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

func Send(subject, content, from, to string) error {
	smtpServer, port, user, password := getSMTPConfig()

	tmpl, err := loadTemplate("templates/templates.html")
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	emailBody, err := generateEmailBody(subject, content, from, tmpl)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	if err := sendMail(smtpServer, port, user, password, from, to, emailBody); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func getSMTPConfig() (server, port, user, password string) {
	return os.Getenv("MAIL_HOST"),
		os.Getenv("MAIL_PORT"),
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD")
}

func loadTemplate(templatePath string) (*template.Template, error) {
	return template.ParseFiles(templatePath)
}

func generateEmailBody(subject, content, from string, tmpl *template.Template) (string, error) {
	data := map[string]string{
		"Subject": subject,
		"Content": content,
		"From":    from,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", subject, body.String()), nil
}

func sendMail(server, port, user, password, from, to, body string) error {
	auth := smtp.PlainAuth("", user, password, server)
	return smtp.SendMail(server+":"+port, auth, from, []string{to}, []byte(body))
}
