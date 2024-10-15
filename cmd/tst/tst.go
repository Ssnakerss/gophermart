package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/logger"
	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/user"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const (
	TokenExp  = time.Hour * 3
	SecretKey = "secretkey"
	//token generate hex
)

func main() {
	logger.Setup("DEV")
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

	bctx := context.Background()

	ddd, err := db.New(db.ConString, db.Info)
	if err != nil {
		log.Fatal(err)
	}

	cr := models.UserCred{
		Login:    "ivan",
		Password: "123456",
	}

	u := user.NewUserManager(ddd)
	// usr, err := u.Register(bctx, &cr)
	// fmt.Println("register:", usr, err)

	usr, err := u.Login(bctx, &cr)
	fmt.Println("login: ", usr, err)

	// // var tokenAuth *jwtauth.JWTAuth

	// tokenAuth := jwtauth.New("HS256", []byte("secret_key_here"), nil)

	// _, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": usr.ID})

	// fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)

}

//type WithContexFunc func(http.ResponseWriter, *http.Request) http.HandlerFunc

// func helloHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Hello, World!")
// }

// func GetUserID(tokenString string) string {
// 	claims := &Claims{}
// 	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
// 		return []byte(SecretKey), nil
// 	})
// 	return claims.UserID
// }

// func BuildJWTString() (string, error) {
// 	token := jwt.NewWithClaims(
// 		jwt.SigningMethodHS256,
// 		Claims{
// 			RegisteredClaims: jwt.RegisteredClaims{
// 				ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
// 			},
// 			UserID: "ivan",
// 		})
// 	tokenString, err := token.SignedString([]byte(SecretKey))
// 	if err != nil {
// 		return "", err
// 	}
// 	return tokenString, nil
// }
