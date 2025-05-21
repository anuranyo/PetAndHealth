package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailConfig struct {
	SMTPHost string
	SMTPPort string
	Username string
	Password string
	FromName string
}

func GetEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
	}
}

func SendPasswordEmail(to, fullName, password, role string) error {
	config := GetEmailConfig()

	subject := "Your PetHealth Account Details"
	body := fmt.Sprintf(`Hello %s,
		An administrator has created a %s account for you in the PetHealth system.

		Your login details:
		Email: %s
		Temporary Password: %s

		Please login to the system and change your password as soon as possible.

		Regards,
		PetHealth Team
	`, fullName, role, to, password)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", config.FromName, config.Username)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=utf-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + body

	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPHost)
	err := smtp.SendMail(
		config.SMTPHost+":"+config.SMTPPort,
		auth,
		config.Username,
		[]string{to},
		[]byte(message),
	)

	return err
}

func SendPasswordResetEmail(to, fullName, tempPassword, role string) error {
	config := GetEmailConfig()

	subject := "PetHealth Password Reset"
	body := fmt.Sprintf(`Hello %s,
		We received a request to reset your password for your %s account in the PetHealth system.

		Your new temporary password is: %s

		Please login using this temporary password and change it immediately for security reasons.

		If you did not request a password reset, please contact support immediately.

		Regards,
		PetHealth Team
	`, fullName, role, tempPassword)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", config.FromName, config.Username)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=utf-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + body

	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPHost)
	err := smtp.SendMail(
		config.SMTPHost+":"+config.SMTPPort,
		auth,
		config.Username,
		[]string{to},
		[]byte(message),
	)

	return err
}
