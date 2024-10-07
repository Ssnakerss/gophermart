package main

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/Ssnakerss/gophermart/internal/logger"
	"github.com/Ssnakerss/gophermart/internal/mock"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const (
	TOKEN_EXP  = time.Hour * 3
	SECRET_KEY = "secretkey"
	//token generate hex
)

func main() {
	logger.Setup("DEV")
	slog.Info("Hello", "module", "tst")

	tokenString, err := BuildJWTString()
	if err != nil {
		log.Fatal((err))
	}

	slog.Info("created", "JWT", tokenString)

	for i := 0; i < 10; i++ {
		_, s, a := mock.YourAcrualIs(i)
		fmt.Printf("order %d status is %s and acrual is %f\n\r", i, s, a)
	}
}

func GetUserId(tokenString string) string {
	claims := &Claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	return claims.UserID
}

func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
			},
			UserID: "ivan",
		})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
