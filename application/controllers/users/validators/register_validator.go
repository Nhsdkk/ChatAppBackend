package validators

import (
	"chat_app_backend/application/models/exception"
	"errors"
	"fmt"
	"regexp"
	"time"
)

const minAcceptableAge = 18

var emailRegexp = regexp.MustCompile("^\\S+@\\w+\\.\\w{2,4}$")
var specialSymbolRegexp = regexp.MustCompile("[\\_\\-\\*\\&\\#\\@\\%]")
var numbersRegexp = regexp.MustCompile("\\d")
var lowercaseLettersRegexp = regexp.MustCompile("[a-z]")
var uppercaseLettersRegexp = regexp.MustCompile("[A-Z]")

func ValidatePassword(password string) error {
	if !numbersRegexp.MatchString(password) ||
		!lowercaseLettersRegexp.MatchString(password) ||
		!uppercaseLettersRegexp.MatchString(password) ||
		!specialSymbolRegexp.MatchString(password) {
		return exception.InvalidBodyException{
			Err: errors.New("password should have at least one of each of this characters (special characters, upper and lowercase letters, digits)"),
		}
	}

	return nil
}

func ValidateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return exception.InvalidBodyException{
			Err: errors.New("email is of wrong format"),
		}
	}
	return nil
}

func ValidateBirthDate(birthDate time.Time) error {
	now := time.Now()
	if birthDate.After(now) {
		return exception.InvalidBodyException{
			Err: errors.New("birth date can't be after today"),
		}
	}

	years := now.Year() - birthDate.Year()
	turningPointThisYear := birthDate.AddDate(years, 0, 0)
	if turningPointThisYear.After(now) {
		years--
	}

	if years < minAcceptableAge {
		return exception.InvalidBodyException{
			Err: errors.New(
				fmt.Sprintf(
					"you are not old enough to register as you should be older than %v to do it",
					minAcceptableAge,
				),
			),
		}
	}

	return nil
}
