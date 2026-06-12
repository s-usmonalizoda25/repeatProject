package repository

import (
	"context"
	"database/sql"
	"project/internal/models"
)

type IUserRepo interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByName(ctx context.Context, name string) (*models.User, error)
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	SoftDelete(ctx context.Context, id int) error
}

type userRepository struct {
	db *sql.DB
}

func New(db *sql.DB) IUserRepo {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	return []models.User{}, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	return &models.User{}, nil
}

func (r *userRepository) GetByName(ctx context.Context, name string) (*models.User, error) {
	return &models.User{}, nil
}

func (r *userRepository) Create(ctx context.Context, user models.User) error {
	return nil
}

func (r *userRepository) Update(ctx context.Context, user models.User) error {
	return nil
}

func (r *userRepository) SoftDelete(ctx context.Context, id int) error {
	return nil
}
