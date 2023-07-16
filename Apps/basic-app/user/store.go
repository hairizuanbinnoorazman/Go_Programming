package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Store interface {
	Create(ctx context.Context, u User) error
	Get(ctx context.Context, ID string) (User, error)
	GetUserByEmail(ctx context.Context, Email string) (User, error)
	GetUserByActivationToken(ctx context.Context, ActivationToken string) (User, error)
	GetUserByForgetPasswordToken(ctx context.Context, ForgetPasswordToken string) (User, error)
	Update(ctx context.Context, ID string, setters ...func(*User) error) (User, error)
}

func setForgetPasswordToken() func(*User) error {
	return func(a *User) error {
		a.ForgetPasswordToken = uuid.New().String()
		a.ForgetPasswordExpiryDate = time.Now().Add(1 * time.Hour)
		return nil
	}
}

// This function is needed in the case we need to resend activation link once more
func setActivationToken() func(*User) error {
	return func(a *User) error {
		a.ActivationToken = uuid.New().String()
		a.ActivationExpiryDate = time.Now().Add(1 * time.Hour)
		return nil
	}
}

func setNewPassword(password string) func(*User) error {
	return func(a *User) error {
		a.ForgetPasswordToken = ""
		a.ForgetPasswordExpiryDate = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		err := a.setPassword(password)
		if err != nil {
			return err
		}
		return nil
	}
}

func setActivate() func(*User) error {
	return func(a *User) error {
		a.ActivationToken = ""
		a.ActivationExpiryDate = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		a.Activated = true
		return nil
	}
}
