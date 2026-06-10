package utils

import (
	"strings"
	"unicode"
)

func ValidateEmail(email string) (string) {
	if email == "" {
		return "Почта должна быть не пустой"
	}
	
	atIndex := strings.Index(email, "@")
	if atIndex < 1 {
		return "Почта должна быть именована и иметь знак @"
	}
	
	dotIndex := strings.Index(email[atIndex:], ".")
	if dotIndex < 2 {
		return "Почта должна быть именована и иметь знак @"
	}
	
	return ""
}

func ValidatePassword(password string) (string) {
	if len(password) < 6 {
		return "Пароль должен содержать минимум 6 символов"
	}
	
	hasDigit := false
	hasUpper := false

	for _, ch := range password {
		if unicode.IsDigit(ch) {
			hasDigit = true
		}
		if unicode.IsUpper(ch) {
			hasUpper = true
		}
	}
	
	if !hasDigit {
		return "Пароль должен содержать хотя бы одну цифру"
	}
	if !hasUpper {
		return "Пароль должен содержать хотя бы одну заглавную букву"
	}
	
	return ""
}

func ValidateRequired(value string) (string) {
	if value == "" {
		return "Данное поле не должно быть пустым"
	}
	return ""
}