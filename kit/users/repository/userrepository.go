package repository

import (
	"context"

	"github.com/juaguz/yuno/kit/users/dto"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) FindByExternalID(ctx context.Context, externalID string) (*dto.User, error) {
	var user dto.User
	err := u.db.WithContext(ctx).Where("user_id = ?", externalID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
