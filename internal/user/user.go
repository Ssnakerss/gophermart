package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/Ssnakerss/gophermart/internal/models"
)

type User models.User

func (u *User) Login(id string, pass string) {
	u.ID = id
	hash, err := makeHash(id, pass)
	if err != nil {
		return
	}
	u.Hash = hash
	u.IsAuthorized = true
}

func (u *User) Register(id string, pass string) {
	hash, err := makeHash(id, pass)
	if err != nil {
		return
	}
	u.Hash = hash
	u.IsAuthorized = true
}

func makeHash(id string, pass string) (string, error) {
	hash := ``
	h := hmac.New(sha256.New, []byte(pass))
	_, err := h.Write([]byte(id))
	if err != nil {
		return ``, err
	}
	hash = hex.EncodeToString(h.Sum(nil))
	return hash, nil
}
