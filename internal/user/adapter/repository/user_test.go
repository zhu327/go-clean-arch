package repository

import (
	"errors"
	"testing"

	"go-clean-arch/internal/user/usecase"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestClassifyCreateErrorMapsEmailUniqueViolationToApplicationConflict(t *testing.T) {
	err := classifyCreateError(&pgconn.PgError{Code: "23505", ConstraintName: "users_email_key"})

	var appErr *usecase.ApplicationError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T, want ApplicationError", err)
	}
	if appErr.Code != usecase.ErrorCodeEmailAlreadyExists || appErr.HTTPStatus != usecase.StatusConflict {
		t.Fatalf("application error = %+v", appErr)
	}
}

func TestClassifyCreateErrorMapsUsernameUniqueViolationToApplicationConflict(t *testing.T) {
	err := classifyCreateError(&pgconn.PgError{Code: "23505", ConstraintName: "users_username_key"})

	var appErr *usecase.ApplicationError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T, want ApplicationError", err)
	}
	if appErr.Code != usecase.ErrorCodeUsernameAlreadyExists || appErr.HTTPStatus != usecase.StatusConflict {
		t.Fatalf("application error = %+v", appErr)
	}
}

func TestClassifyCreateErrorLeavesOtherErrorsUntouched(t *testing.T) {
	original := errors.New("database unavailable")
	if got := classifyCreateError(original); !errors.Is(got, original) {
		t.Fatalf("error = %v, want wrapped original %v", got, original)
	}
}
