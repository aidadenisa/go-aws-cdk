package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(registerUser RegisterUser) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		return User{}, err
	}

	return User{
		Username:registerUser.Username,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ValidatePassword(hashedPasswork, plainTextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswork), []byte(plainTextPassword))
	return err == nil
}

func CreateToken(user User) string {
	now := time.Now()
	
	validUntil := now.Add(time.Hour).Unix()

	// JWT payload
	claims := jwt.MapClaims{
		"user": user.Username,
		"expires": validUntil,
	}

	// Most popular method: jwt.SigningMethodHS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := "secret"

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return tokenString

}