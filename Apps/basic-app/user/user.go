package user

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	errEmailInvalid               = errors.New("email is invalid")
	errPasswordShort              = errors.New("password cannot be shorter than 8 characters")
	errPasswordLong               = errors.New("password cannot be longer than 120 characters")
	errPasswordInvalid            = errors.New("password requires at least 1 capital letter, 1 small letter and a number")
	errSamePassword               = errors.New("current password already in use. Please pick another password")
	errActivationTokenInvalid     = errors.New("activation Token is invalid")
	errForgetPasswordTokenInvalid = errors.New("forget Password Token is invalid")
)

type User struct {
	ID                       string    `json:"id" gorm:"type:varchar(40);primary_key"`
	Email                    string    `json:"email" gorm:"type:varchar(250)"`
	Password                 string    `json:"-" gorm:"type:varchar(250)"`
	ForgetPasswordToken      string    `json:"forget_password_token" gorm:"type:varchar(40)"`
	ForgetPasswordExpiryDate time.Time `json:"forget_password_expiry_date"`
	ActivationToken          string    `json:"activation_token" gorm:"type:varchar(40)"`
	ActivationExpiryDate     time.Time `json:"activation_expiry_date"`
	Activated                bool      `json:"activated"`
	DateCreated              time.Time `json:"date_created"`
	DateModified             time.Time `json:"date_modified"`
}

func New(email, password string) (User, error) {
	user := User{Email: email}
	err := user.setPassword(password)
	user.ID = uuid.New().String()
	user.ActivationToken = uuid.New().String()
	user.DateCreated = time.Now()
	user.DateModified = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	user.ForgetPasswordExpiryDate = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	user.ActivationExpiryDate = time.Now().Add(1 * time.Hour)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *User) setPassword(password string) error {
	if len(password) < 8 {
		return errPasswordShort
	}
	if len(password) > 120 {
		return errPasswordLong
	}
	reSmallLetters := regexp.MustCompile("[a-z]")
	reCapital := regexp.MustCompile("[A-Z]")
	reNumbers := regexp.MustCompile("[0-9]")
	smallLettersFind := reSmallLetters.FindAllString(password, -1)
	capitalFind := reCapital.FindAllString(password, -1)
	numberFind := reNumbers.FindAllString(password, -1)
	if len(smallLettersFind) > 0 && len(capitalFind) > 0 && len(numberFind) > 0 {
		hashedPassword, errEncrpt := bcrypt.GenerateFromPassword([]byte(password), 10)
		if errEncrpt != nil {
			return errPasswordInvalid
		}
		u.Password = string(hashedPassword)
		return nil
	}
	return errPasswordInvalid
}
