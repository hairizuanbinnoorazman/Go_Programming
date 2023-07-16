package user

import (
	"context"
	"errors"

	"github.com/hairizuanbinnoorazman/basic-app/logger"
	"github.com/jinzhu/gorm"
)

type mysql struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewMySQL(logger logger.Logger, dbClient *gorm.DB) mysql {
	return mysql{
		db:     dbClient,
		logger: logger,
	}
}

func (m mysql) Create(ctx context.Context, u User) error {
	result := m.db.Create(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m mysql) GetUser(ctx context.Context, ID string) (User, error) {
	u := User{}
	result := m.db.Where("id = ?", ID).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, nil
	}
	if result.Error != nil {
		return User{}, result.Error
	}
	return u, nil
}

func (m mysql) GetUserByEmail(ctx context.Context, Email string) (User, error) {
	u := User{}
	result := m.db.Where("email = ?", Email).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, nil
	}
	if result.Error != nil {
		return User{}, result.Error
	}
	return u, nil
}

func (m mysql) GetUserByActivationToken(ctx context.Context, ActivationToken string) (User, error) {
	u := User{}
	result := m.db.Where("activation_token = ?", ActivationToken).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, nil
	}
	if result.Error != nil {
		return User{}, result.Error
	}
	return u, nil
}

func (m mysql) GetUserByForgetPasswordToken(ctx context.Context, ForgetPasswordToken string) (User, error) {
	u := User{}
	result := m.db.Where("forget_password_token = ?", ForgetPasswordToken).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, nil
	}
	if result.Error != nil {
		return User{}, result.Error
	}
	return u, nil
}

func (m mysql) Update(ctx context.Context, ID string, setters ...func(*User) error) (User, error) {
	var u User
	result := m.db.Where("id = ?", ID).First(&u)
	if result.Error != nil {
		return User{}, result.Error
	}
	for _, s := range setters {
		err := s(&u)
		if err != nil {
			return User{}, err
		}
	}
	result = m.db.Save(&u)
	if result.Error != nil {
		return User{}, result.Error
	}

	return u, nil
}
