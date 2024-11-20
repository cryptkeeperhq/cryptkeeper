package utils

import (
	"math/rand"
)

const (
	lowerCharset   = "abcdefghijklmnopqrstuvwxyz"
	upperCharset   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitCharset   = "0123456789"
	specialCharset = "!@#$%^&*()-_=+[]{}|;:,.<>?/~`"
	allCharset     = lowerCharset + upperCharset + digitCharset + specialCharset
	passwordLength = 16
)

const usernameCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const usernameLength = 12

func GenerateUsername() string {
	b := make([]byte, usernameLength)
	for i := range b {
		b[i] = usernameCharset[rand.Intn(len(usernameCharset))]
	}
	return string(b)
}

func GeneratePassword() string {
	b := make([]byte, passwordLength)
	for i := range b {
		b[i] = allCharset[rand.Intn(len(allCharset))]
	}
	return string(b)
}
