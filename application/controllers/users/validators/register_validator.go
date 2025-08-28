package validators

import (
	"regexp"
	"time"
)

const minAcceptableAge = 18

var emailRegexp = regexp.MustCompile("^\\S+@\\w+\\.\\w{2,4}$")
var specialSymbolRegexp = regexp.MustCompile("[\\_\\-\\*\\&\\#\\@\\%]")
var numbersRegexp = regexp.MustCompile("\\d")
var lowercaseLettersRegexp = regexp.MustCompile("[a-z]")
var uppercaseLettersRegexp = regexp.MustCompile("[A-Z]")

type PasswordValidator struct{}

func (p PasswordValidator) Validate(password *string) bool {
	if !numbersRegexp.MatchString(*password) ||
		!lowercaseLettersRegexp.MatchString(*password) ||
		!uppercaseLettersRegexp.MatchString(*password) ||
		!specialSymbolRegexp.MatchString(*password) {
		return false
	}

	return true
}

type EmailValidator struct{}

func (e EmailValidator) Validate(email *string) bool {
	if !emailRegexp.MatchString(*email) {
		return false
	}

	return true
}

type BirthDateValidator struct{}

func (b BirthDateValidator) Validate(birthDate *time.Time) bool {
	now := time.Now()
	if birthDate.After(now) {
		return false
	}

	years := now.Year() - birthDate.Year()
	turningPointThisYear := birthDate.AddDate(years, 0, 0)
	if turningPointThisYear.After(now) {
		years--
	}

	if years < minAcceptableAge {
		return false
	}

	return true
}
