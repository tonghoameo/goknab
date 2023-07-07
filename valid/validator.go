package valid

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(s string, minLen int, maxLen int) error {
	n := len(s)
	if n < minLen || n > maxLen {
		return fmt.Errorf("must contain from %d to %d \n", minLen, maxLen)
	}

	return nil
}
func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 10); err != nil {
		return err
	}
	if !isValidUsername(username) {
		return fmt.Errorf("must contain only lower letter and number and underscore [%s]", username)
	}
	return nil
}
func ValidateFullname(fullname string) error {
	if err := ValidateString(fullname, 3, 10); err != nil {
		return err
	}
	if !isValidFullname(fullname) {
		return fmt.Errorf("Fullname must contain only letter and number and underscore ")
	}
	return nil
}
func ValidatePassword(password string) error {
	return ValidateString(password, 6, 200)

}
func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("is not a valid Email address")
	}

	return nil
}
func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("id cannot negative")
	}
	return nil
}
func ValidateSecretCode(secretCode string) error {
	if err := ValidateString(secretCode, 32, 128); err != nil {
		return err
	}
	return nil
}
