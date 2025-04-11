package users

import (
	"fmt"
	"net/mail"
	"strings"
)

func ValidateEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	if _, err := mail.ParseAddress(email); err != nil {
		return "", err
	}
	return email, nil
}

func ValidatePassword(password string) (string, error) {
	password = strings.TrimSpace(password)
	if len(password) <= 8 {
		return "", fmt.Errorf("password too short")
	}
	return password, nil
}
