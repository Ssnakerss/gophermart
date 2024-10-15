package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func (hm *HandlerMaster) AuthorizeUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		fmt.Printf("claims %v\r\n", claims)

		if len(claims) == 0 {
			w.WriteHeader(http.StatusUnauthorized) // 401
			return
		}
		hm.currentUser.ID = claims["user_id"].(string)
		err := hm.UserManager.CheckUserExist(hm.rootAppContext, hm.currentUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized) // 401
			return
		}

		next.ServeHTTP(w, r) // continue
	})
}
