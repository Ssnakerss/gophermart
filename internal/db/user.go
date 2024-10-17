package db

import (
	"context"
	"errors"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func (db *GormDB) CreateUser(ctx context.Context, user *models.User) error {
	err := db.DB.WithContext(ctx).Create(user).Error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		//если пользователь уже существует, возвращаем ошибку
		if pgErr.Code == "23505" {
			return models.ErrUserAlreadyExists
		}
	}
	return err
}

func (db *GormDB) GetUser(ctx context.Context, user *models.User) error {
	err := db.DB.
		WithContext(ctx).
		Where("hash = ?", user.Hash).
		First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ErrUserNotFound
	}
	return err
}
func (db *GormDB) CheckUserExist(ctx context.Context, user *models.User) error {
	err := db.DB.WithContext(ctx).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ErrUserNotFound
	}
	return err
}
