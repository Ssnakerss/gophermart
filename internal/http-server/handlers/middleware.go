package handlers

import (
	"net/http"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/go-chi/jwtauth/v5"
)

func (hm *HandlerMaster) AuthorizeUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		if len(claims) == 0 {
			w.WriteHeader(http.StatusUnauthorized) // 401
			return
		}

		hm.currentUser = &models.User{ID: claims["user_id"].(string)}
		err := hm.UserManager.CheckUserExist(hm.rootAppContext, hm.currentUser)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized) // 401
			return
		}

		next.ServeHTTP(w, r) // continue
	})
}
