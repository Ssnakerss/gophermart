package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/go-chi/jwtauth/v5"
)

type UserManager struct {
	storage   models.UserStorage
	tokenAuth *jwtauth.JWTAuth
}

func NewUserManager(storage models.UserStorage) *UserManager {
	return &UserManager{
		storage:   storage,
		tokenAuth: jwtauth.New("HS256", []byte("secret_key_here!"), nil), //TODO - сделать секрет ки
	}
}

type UserID string

func (u *UserManager) Login(ctx context.Context, cred *models.UserCred) (*models.User, error) {
	var err error
	user := &models.User{ID: cred.Login}
	user.Hash, err = makeHash(cred)
	if err != nil {
		return user, err
	}
	err = u.storage.GetUser(ctx, user)
	return user, err
}

func (u *UserManager) Register(ctx context.Context, cred *models.UserCred) (*models.User, error) {
	var err error
	user := &models.User{ID: cred.Login}
	user.Hash, err = makeHash(cred)
	if err != nil {
		return user, err
	}
	//создаем пользователя
	err = u.storage.CreateUser(ctx, user)

	return user, err
}

func (u *UserManager) CreateJWT(user *models.User) (string, error) {
	_, tokenString, err := u.tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
	return tokenString, err
}

func (u *UserManager) CheckUserExist(ctx context.Context, user *models.User) error {
	return u.storage.CheckUserExist(ctx, user)
}

func (u *UserManager) GetTokenAuth() *jwtauth.JWTAuth {
	return u.tokenAuth
}

func makeHash(cred *models.UserCred) (string, error) {
	hash := ``
	h := hmac.New(sha256.New, []byte(cred.Password))
	_, err := h.Write([]byte(cred.Login))
	if err != nil {
		return ``, err
	}
	hash = hex.EncodeToString(h.Sum(nil))
	return hash, nil
}
