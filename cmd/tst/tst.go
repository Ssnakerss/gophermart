package main

import (
	"fmt"
	"time"

	"github.com/Ssnakerss/gophermart/internal/models"
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

type money uint64

func main() {
	// logger.Setup("DEV")
	// slog.Info("Hello", "module", "tst")

	// tokenString, err := BuildJWTString()
	// if err != nil {
	// 	log.Fatal((err))
	// }

	// slog.Info("created", "JWT", tokenString)

	// for i := 0; i < 10; i++ {
	// 	_, s, a := mock.YourAcrualIs(i)
	// 	fmt.Printf("order %d status is %s and acrual is %f\n\r", i, s, a)
	// }

	// var sum float64
	// for i := 0; i < 10_000_000; i++ {
	// 	sum += float64(0.8)
	// }

	// var expectedSum int

	// for i := 0; i < 1000_000; i++ {
	// 	expectedSum += int(8)
	// }

	// fmt.Println(sum, expectedSum, sum == float64(expectedSum))

	var b models.Bonus

	b.Set(1.11)
	fmt.Printf("bonus: %d, currency: %f\n\r", b, b.Get())
	b = b + 100
	fmt.Printf("new bonus: %d, currency: %f\n\r", b, b.Get())
	b.Add(100.33)
	fmt.Printf("after add bonus: %d, currency: %f\n\r", b, b.Get())
	b.Sub(100.11)
	fmt.Printf("after sub bonus: %d, currency: %f\n\r", b, b.Get())

	fmt.Printf("trying to sub more %s \n\r", b.Sub(1000))
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
