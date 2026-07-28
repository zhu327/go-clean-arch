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

func TestNewUser_RejectsShortUsername(t *testing.T) {
	_, err := NewUser("ab", "test@test.com", "hashed")
	if !errors.Is(err, ErrUsernameTooShort) {
		t.Fatalf("err = %v, want ErrUsernameTooShort", err)
	}
}

func TestNewUser_RejectsLongUsername(t *testing.T) {
	_, err := NewUser("abcdefghijklmnopqrstuvwxyz", "test@test.com", "hashed")
	if !errors.Is(err, ErrUsernameTooLong) {
		t.Fatalf("err = %v, want ErrUsernameTooLong", err)
	}
}

func TestNewUser_RejectsInvalidEmail(t *testing.T) {
	for _, email := range []string{"not-an-email", "missing@tld", "@test.com", "test@.com"} {
		_, err := NewUser("testuser", email, "hashed")
		if !errors.Is(err, ErrInvalidEmail) {
			t.Fatalf("email %q: err = %v, want ErrInvalidEmail", email, err)
		}
	}
}

func TestNewUser_RejectsEmptyPassword(t *testing.T) {
	_, err := NewUser("testuser", "test@test.com", "")
	if !errors.Is(err, ErrEmptyPassword) {
		t.Fatalf("err = %v, want ErrEmptyPassword", err)
	}
}

func TestNewUser_AcceptsValidInput(t *testing.T) {
	user, err := NewUser("testuser", "test@test.com", "hashed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "testuser" || user.Email != "test@test.com" || user.Password != "hashed" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestChangePassword_RejectsEmptyPassword(t *testing.T) {
	user := &User{}
	if err := user.ChangePassword(""); !errors.Is(err, ErrEmptyPassword) {
		t.Fatalf("err = %v, want ErrEmptyPassword", err)
	}
}

func TestChangePassword_UpdatesPassword(t *testing.T) {
	user := &User{Password: "old"}
	if err := user.ChangePassword("new-hashed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Password != "new-hashed" {
		t.Errorf("Password = %q, want %q", user.Password, "new-hashed")
	}
}
