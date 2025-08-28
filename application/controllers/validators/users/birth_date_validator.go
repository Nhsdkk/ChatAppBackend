package user_validators

import (
	"context"
	"time"
)

const minAcceptableAge = 18

type BirthDateValidator struct{}

func (b BirthDateValidator) Validate(birthDate *time.Time, _ context.Context) bool {
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
