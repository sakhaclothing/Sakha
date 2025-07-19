package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/sakhaclothing/config"
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
	smtpConfig, err := config.GetSMTPConfig()
	if err != nil {
		return fmt.Errorf("SMTP config not found: %v", err)
	}

	// If SMTP credentials are not configured, use mock email
	if smtpConfig.SMTPUsername == "" || smtpConfig.SMTPPassword == "" {
		return SendMockPasswordResetEmail(toEmail, resetToken, resetLink)
	}

	subject := "Reset Password - Sakha Clothing"
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
			<div style="background-color: #f8f9fa; padding: 30px; border-radius: 10px;">
				<h2 style="color: #333; text-align: center; margin-bottom: 30px;">Reset Password</h2>
				
				<p style="color: #666; line-height: 1.6;">Halo,</p>
				
				<p style="color: #666; line-height: 1.6;">Anda telah meminta untuk mereset password akun Sakha Clothing Anda.</p>
				
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #007bff; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
				</div>
				
				<p style="color: #666; line-height: 1.6;">Atau copy paste link berikut ke browser Anda:</p>
				<p style="background-color: #f1f3f4; padding: 10px; border-radius: 5px; word-break: break-all; color: #333;">%s</p>
				
				<div style="background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 5px; margin: 20px 0;">
					<p style="color: #856404; margin: 0; font-size: 14px;">
						<strong>Penting:</strong> Link ini akan expired dalam 1 jam. Jika Anda tidak meminta reset password, abaikan email ini.
					</p>
				</div>
				
				<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
				
				<p style="color: #999; font-size: 12px; text-align: center;">
					Best regards,<br>
					<strong>Sakha Clothing Team</strong>
				</p>
			</div>
		</body>
		</html>
	`, resetLink, resetLink)

	return SendSMTPEmail(smtpConfig, toEmail, subject, body)
}

// SendSMTPEmail sends email via SMTP
func SendSMTPEmail(smtpConfig *config.SMTPConfig, toEmail, subject, body string) error {
	auth := smtp.PlainAuth("", smtpConfig.SMTPUsername, smtpConfig.SMTPPassword, smtpConfig.SMTPHost)

	// Email headers
	headers := fmt.Sprintf("From: %s <%s>\r\n", smtpConfig.FromName, smtpConfig.FromEmail)
	headers += fmt.Sprintf("To: %s\r\n", toEmail)
	headers += fmt.Sprintf("Subject: %s\r\n", subject)
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=UTF-8\r\n"
	headers += "\r\n"

	message := headers + body

	// Send email
	err := smtp.SendMail(
		smtpConfig.SMTPHost+":"+smtpConfig.SMTPPort,
		auth,
		smtpConfig.FromEmail,
		[]string{toEmail},
		[]byte(message),
	)

	if err != nil {
		fmt.Printf("SMTP Error: %v\n", err)
		return err
	}

	fmt.Printf("Email sent successfully to: %s\n", toEmail)
	return nil
}

// SendEmail is a generic email sending function
func SendEmail(toEmail, subject, body string) error {
	smtpConfig, err := config.GetSMTPConfig()
	if err != nil {
		return fmt.Errorf("SMTP config not found: %v", err)
	}

	// If SMTP credentials are not configured, use mock email
	if smtpConfig.SMTPUsername == "" || smtpConfig.SMTPPassword == "" {
		return SendMockEmail(toEmail, subject, body)
	}

	return SendSMTPEmail(smtpConfig, toEmail, subject, body)
}

// SendMockEmail is a mock function for development
func SendMockEmail(toEmail, subject, body string) error {
	// In development, just log the email details
	fmt.Printf("=== MOCK EMAIL (SMTP not configured) ===\n")
	fmt.Printf("To: %s\n", toEmail)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("Body: %s\n", body)
	fmt.Printf("========================================\n")

	return nil
}

// SendMockPasswordResetEmail is a mock function for development
func SendMockPasswordResetEmail(toEmail, resetToken, resetLink string) error {
	// In development, just log the email details
	fmt.Printf("=== MOCK EMAIL (SMTP not configured) ===\n")
	fmt.Printf("To: %s\n", toEmail)
	fmt.Printf("Subject: Reset Password - Sakha Clothing\n")
	fmt.Printf("Reset Token: %s\n", resetToken)
	fmt.Printf("Reset Link: %s\n", resetLink)
	fmt.Printf("========================================\n")

	return nil
}

// SendTLSEmail sends email with TLS (alternative method)
func SendTLSEmail(smtpConfig *config.SMTPConfig, toEmail, subject, body string) error {
	// Create TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpConfig.SMTPHost,
	}

	// Connect to SMTP server
	conn, err := tls.Dial("tcp", smtpConfig.SMTPHost+":"+smtpConfig.SMTPPort, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, smtpConfig.SMTPHost)
	if err != nil {
		return err
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", smtpConfig.SMTPUsername, smtpConfig.SMTPPassword, smtpConfig.SMTPHost)
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Set sender
	if err = client.Mail(smtpConfig.FromEmail); err != nil {
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
	headers := fmt.Sprintf("From: %s <%s>\r\n", smtpConfig.FromName, smtpConfig.FromEmail)
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
