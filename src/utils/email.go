package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, otp string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	htmlBody := fmt.Sprintf(`
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; color: #333; }
				.container { padding: 20px; text-align: center; }
				.otp { font-size: 36px; font-weight: bold; color: #333; letter-spacing: 10px; }
				.note { font-size: 14px; margin-top: 20px; color: #555; }
			</style>
		</head>
		<body>
			<div class="container">
				<p>To reset your password, please use the following OTP code:</p>
				<div class="otp">%s</div>
				<p class="note">
					Please do not share this OTP to anyone, including people claiming from official website.<br><br>
					This code will expire in 5 minutes.<br><br>
					If you did not request this code, please ignore this message.
				</p>
			</div>
		</body>
		</html>`, otp)

	// Format email
	msg := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		htmlBody

	// Setup SMTP auth
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Kirim email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
