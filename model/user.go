package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
)

const passwordMinLength = 8
const passwordMaxLength = 256

// User represents the user of the service.
// First the user is registered in the GopherMart loyality system. When authenticated
// user makes a purchase in GopherMart store, the order (with information about the goods and
// their prices) is sent to the GopherPoints service. Also the order number is sent to the UserService
// and that service makes a request to the GopherPoint service with a user ID and order number and
// receives a number of bonus G-Points. And those points are added to user's bonus balance.
type User struct {
	ID    uuid.UUID
	Login string
	// Password value is deleted after encrypting.
	Password string
	// PasswordHash is bcrypt hashed user's password.
	PasswordHash  string
	CreatedAt     time.Time
	RememberToken string
	// GPointsBalance is user's bonus account balance
	GPointsBalance float32
}

// Validate performs User fields checking.
func (c User) Validate() error {
	var result *multierror.Error

	if err := validateLogin(c.Login); err != nil {
		result = multierror.Append(result, err)
	}
	if err := validatePassword(c.Password); err != nil {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}

func validateLogin(login string) error {
	if len(login) < 3 || len(login) > 64 {
		return errors.New("validate login: invalid length")
	}

	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return errors.New("valdate password: empty")
	}
	if len(password) < passwordMinLength || len(password) > passwordMaxLength {
		return errors.New("validate password: invalid length")
	}

	// TODO: we can check if there is at least one digit, capital letter etc.

	return nil
}
