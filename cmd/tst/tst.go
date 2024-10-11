package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/logger"
	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	// var sum float64
	// for i := 0; i < 10_000_000; i++ {
	// 	sum += float64(0.8)
	// }

	// var expectedSum int

	// for i := 0; i < 1000_000; i++ {
	// 	expectedSum += int(8)
	// }

	// fmt.Println(sum, expectedSum, sum == float64(expectedSum))

	// var b types.Bonus

	// b.Set(1.11)
	// fmt.Printf("bonus: %d, currency: %f\n\r", b, b.Get())
	// b = b + 100
	// fmt.Printf("new bonus: %d, currency: %f\n\r", b, b.Get())
	// b.Add(100.33)
	// fmt.Printf("after add bonus: %d, currency: %f\n\r", b, b.Get())
	// b.Sub(100.11)
	// fmt.Printf("after sub bonus: %d, currency: %f\n\r", b, b.Get())

	// fmt.Printf("trying to sub more %s \n\r", b.Sub(1000))

	// t := time.Now()

	// s := t.Format("2006-01-02T15:04:05Z07:00")
	//                2024-10-10T13:01:17.617559+03:00
	// tt, _ := time.Parse("2006-01-02T15:04:05Z07:00", s)

	// fmt.Printf("RFC3339 : %s | Date : %v", s, tt)
	// bctx := context.Background()

	// // ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	db := db.New(db.ConString, db.Info)
	// db.Migrate(ctx)

	account := models.Account{UserID: "ivan"}
	fmt.Printf("%v\r\n", account)
	r := db.DB.Model(&account).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "balance"}}}).
		Where("balance >= ?", 1).
		Update("balance", gorm.Expr("balance - ?", 1))
	fmt.Printf("%v\r\n", account)

	r = db.DB.Model(&account).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "balance"}}}).
		Where("balance >= ?", 10).
		Update("balance", gorm.Expr("balance - ?", 10))
	fmt.Printf("%v\r\n", account)

	fmt.Println(r.Error)

	// var o models.Order
	// o.Number.Set(0)
	// err := db.GetOrder(ctx, &o)
	// fmt.Printf("get order: %v\n\r", o)

	// if errors.Is(err, models.ErrRecordNotFound) {
	// 	fmt.Println("not found error")
	// }

	// uo := models.Order{
	// 	UserID: "dummy",
	// }
	// orders := db.GetAllOrders(ctx, &uo)
	// fmt.Printf("Get Orders \r\n%v\n\r", orders)

	// db.SaveOrder(ctx, &order)

	// order.Accrual.Add(10.11)

	// db.UpdateOrder(ctx, &order)

	// var o models.Order
	// o.Number.Set(397_471)
	// db.GetOrder(ctx, &o)
	// fmt.Println(o)
	// o.Status = models.CHECKING
	// db.UpdateOrder(ctx, &o)

	// server := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: http.HandlerFunc(helloHandler),
	// }

	// server.ConnContext = func(ctx context.Context, c net.Conn) context.Context {
	// 	return context.Background()
	// }

}

//type WithContexFunc func(http.ResponseWriter, *http.Request) http.HandlerFunc

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func GetUserID(tokenString string) string {
	claims := &Claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	return claims.UserID
}

func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
			},
			UserID: "ivan",
		})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
