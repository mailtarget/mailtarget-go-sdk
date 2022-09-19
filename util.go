package layang

import (
	"net/mail"
	"regexp"
)

func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsHTML(html string) bool {
	return regexp.MustCompile(`<[/]?([a-zA-Z]+).*?>`).MatchString(html)
}
