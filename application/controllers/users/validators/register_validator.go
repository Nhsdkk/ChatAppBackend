package validators

import (
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
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
		message := "password should have at least one of each of this characters (special characters, upper and lowercase letters, digits)"
		return common_exceptions.InvalidBodyException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
				Message:             message,
			},
		}
	}

	return nil
}

func ValidateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		message := "email is of wrong format"
		return common_exceptions.InvalidBodyException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
				Message:             message,
			},
		}
	}
	return nil
}

func ValidateBirthDate(birthDate time.Time) error {
	now := time.Now()
	if birthDate.After(now) {
		message := "birth date can't be after today"
		return common_exceptions.InvalidBodyException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
				Message:             message,
			},
		}
	}

	years := now.Year() - birthDate.Year()
	turningPointThisYear := birthDate.AddDate(years, 0, 0)
	if turningPointThisYear.After(now) {
		years--
	}

	if years < minAcceptableAge {
		message := fmt.Sprintf(
			"you are not old enough to register as you should be older than %v to do it",
			minAcceptableAge,
		)

		return common_exceptions.InvalidBodyException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
				Message:             message,
			},
		}
	}

	return nil
}
