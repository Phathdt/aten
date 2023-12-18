package storage

import (
	"aten/module/models"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type sqlStorage struct {
	db *gorm.DB
}

func (s *sqlStorage) CreateUser(ctx context.Context, data *models.UserCreate) error {
	db := s.db.WithContext(ctx).Table(models.User{}.TableName())
	if err := db.Create(&data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *sqlStorage) GetUserByCondition(ctx context.Context, cond map[string]interface{}) (*models.User, error) {
	var data models.User
	db := s.db.WithContext(ctx).Table(models.User{}.TableName())

	result := db.Where(cond).Limit(1).Find(&data)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errors.WithStack(errors.New("not found"))
	}

	return &data, nil
}

func NewSqlStorage(db *gorm.DB) *sqlStorage {
	return &sqlStorage{db: db}
}
