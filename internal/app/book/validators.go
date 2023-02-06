package book

import (
	"log"
	"regexp"
)

var (
	USERNAME_PATTERN *regexp.Regexp
	PHONE_PATTERN *regexp.Regexp
)

func init() {
	var err error
	
	USERNAME_PATTERN, err = regexp.Compile(`^([A-zА-я\s\.]){1,16}$`)
	if err != nil {
		log.Fatal(err)
	}

	PHONE_PATTERN, err = regexp.Compile(`^[+]?\d{1,4}\s?\(?\d{1,3}\)?\s?\d{3}[\s-]?\d{2}[\s-]?\d{2}$`)
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateUsername(u string) bool {
	return USERNAME_PATTERN.MatchString(u)
}

func ValidatePhone(p string) bool {
	return PHONE_PATTERN.MatchString(p)
}

func ValidateContact(c *Contact) bool {
	return ValidateUsername(c.Username) && ValidatePhone(c.Phone)
}