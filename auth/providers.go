package auth

import (
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

type SmsCodeProvider interface {
	GetCode(phone string) (string, error)
}

type ConsoleSmsCodeProvider struct{}

func (c *ConsoleSmsCodeProvider) GetCode(phone string) (string, error) {
	fmt.Printf("Enter SMS code for %s: ", phone)
	var code string
	_, err := fmt.Scanln(&code)
	return strings.TrimSpace(code), err
}

type PasswordProvider interface {
	GetPassword(hint string) (string, error)
}

type ConsolePasswordProvider struct{}

func (c *ConsolePasswordProvider) GetPassword(hint string) (string, error) {
	prompt := "Enter 2FA password"
	if hint != "" {
		prompt += fmt.Sprintf(" (hint: %s)", hint)
	}
	prompt += ": "
	fmt.Print(prompt)
	var password string
	_, err := fmt.Scanln(&password)
	return strings.TrimSpace(password), err
}

type QrHandler interface {
	ShowQR(qrURL string) error
}

type ConsoleQrHandler struct{}

func (c *ConsoleQrHandler) ShowQR(qrURL string) error {
	qr, err := qrcode.New(qrURL, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("failed to generate QR: %w", err)
	}
	fmt.Println(qr.ToSmallString(false))
	fmt.Printf("Or open: %s\n", qrURL)
	return nil
}

type EmailCodeProvider interface {
	GetCode(email string) (string, error)
}
