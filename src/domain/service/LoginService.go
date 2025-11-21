package service

import (
	"errors"
	"os"
	"time"

	"github.com/anaclaraddias/feedmaker-students/src/domain/entity"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const invalidLoginError = "incorrect username or password"

type LoginService struct {
	db *gorm.DB
}

func NewLoginService(db *gorm.DB) *LoginService {
	return &LoginService{
		db: db,
	}
}

func (service *LoginService) Execute(loginUser entity.User) (string, error) {
	var user *entity.User
	if err := service.db.Where("username = ?", loginUser.Username).Find(&user).Error; err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New(invalidLoginError)
	}

	if user.Password != loginUser.Password {
		return "", errors.New(invalidLoginError)
	}

	jwt, err := service.createJwt(user)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (service *LoginService) createJwt(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"id":   user.ID,
		"name": user.Name,
		"type": user.Type,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return jwt, nil
}
