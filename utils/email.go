package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// GetEmailConfig returns email configuration from environment variables
func GetEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@sakha.com"),
		FromName:     getEnv("FROM_NAME", "Sakha Clothing"),
	}
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(toEmail, resetToken, resetLink string) error {
	config := GetEmailConfig()

	// If SMTP credentials are not configured, use mock email
	if config.SMTPUsername == "" || config.SMTPPassword == "" {
		return SendMockPasswordResetEmail(toEmail, resetToken, resetLink)
	}

	subject := "Reset Password - Sakha Clothing"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Reset Password</h2>
			<p>Anda telah meminta untuk mereset password akun Anda.</p>
			<p>Klik link di bawah ini untuk mereset password:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>Atau copy paste link berikut ke browser:</p>
			<p>%s</p>
			<p>Link ini akan expired dalam 1 jam.</p>
			<p>Jika Anda tidak meminta reset password, abaikan email ini.</p>
			<br>
			<p>Best regards,<br>Sakha Clothing Team</p>
		</body>
		</html>
	`, resetLink, resetLink)

	return SendSMTPEmail(config, toEmail, subject, body)
}

// SendSMTPEmail sends email via SMTP
func SendSMTPEmail(config EmailConfig, toEmail, subject, body string) error {
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)

	// Email headers
	headers := fmt.Sprintf("From: %s <%s>\r\n", config.FromName, config.FromEmail)
	headers += fmt.Sprintf("To: %s\r\n", toEmail)
	headers += fmt.Sprintf("Subject: %s\r\n", subject)
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=UTF-8\r\n"
	headers += "\r\n"

	message := headers + body

	// Send email
	err := smtp.SendMail(
		config.SMTPHost+":"+config.SMTPPort,
		auth,
		config.FromEmail,
		[]string{toEmail},
		[]byte(message),
	)

	return err
}

// SendMockPasswordResetEmail is a mock function for development
func SendMockPasswordResetEmail(toEmail, resetToken, resetLink string) error {
	// In development, just log the email details
	fmt.Printf("=== MOCK EMAIL ===\n")
	fmt.Printf("To: %s\n", toEmail)
	fmt.Printf("Subject: Reset Password - Sakha Clothing\n")
	fmt.Printf("Reset Token: %s\n", resetToken)
	fmt.Printf("Reset Link: %s\n", resetLink)
	fmt.Printf("==================\n")

	return nil
}

// SendTLSEmail sends email with TLS (alternative method)
func SendTLSEmail(config EmailConfig, toEmail, subject, body string) error {
	// Create TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.SMTPHost,
	}

	// Connect to SMTP server
	conn, err := tls.Dial("tcp", config.SMTPHost+":"+config.SMTPPort, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, config.SMTPHost)
	if err != nil {
		return err
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Set sender
	if err = client.Mail(config.FromEmail); err != nil {
		return err
	}

	// Set recipient
	if err = client.Rcpt(toEmail); err != nil {
		return err
	}

	// Send email data
	writer, err := client.Data()
	if err != nil {
		return err
	}

	// Email headers
	headers := fmt.Sprintf("From: %s <%s>\r\n", config.FromName, config.FromEmail)
	headers += fmt.Sprintf("To: %s\r\n", toEmail)
	headers += fmt.Sprintf("Subject: %s\r\n", subject)
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=UTF-8\r\n"
	headers += "\r\n"

	message := headers + body

	_, err = writer.Write([]byte(message))
	if err != nil {
		return err
	}

	return writer.Close()
}
