package domain

import (
	"errors"
	"testing"
)

func TestDomainSentinelErrors(t *testing.T) {
	errs := map[string]error{
		"ErrUserNotFound":      ErrUserNotFound,
		"ErrUserAlreadyExists": ErrUserAlreadyExists,
		"ErrInvalidEmail":      ErrInvalidEmail,
		"ErrUsernameTooShort":  ErrUsernameTooShort,
		"ErrUsernameTooLong":   ErrUsernameTooLong,
		"ErrEmptyPassword":     ErrEmptyPassword,
	}
	for name, err := range errs {
		if err == nil {
			t.Errorf("%s is nil", name)
		}
		if !errors.Is(err, err) {
			t.Errorf("%s does not satisfy errors.Is with itself", name)
		}
	}
}
