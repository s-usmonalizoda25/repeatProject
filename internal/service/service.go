package service

import (
	"context"
	"errors"
	"fmt"

	"project/internal/models"
	"project/internal/repository"
	"project/pkg/errs"
)

type UserService struct {
	repo repository.IUserRepo
}

func New(repo repository.IUserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CleanGetAll(ctx context.Context) ([]models.User, error) {
	allUsers, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetAll: %w", err)
	}

	var activeUsers []models.User
	for _, u := range allUsers {
		if u.IsActive {
			activeUsers = append(activeUsers, u)
		}
	}

	return activeUsers, nil
}

func (s *UserService) GetAllArchive(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetAll archive: %w", err)
	}
	return users, nil
}

func (s *UserService) Login(ctx context.Context, name, password string) error {
	if name == "" || password == "" {
		return errs.ErrValidation
	}
	
	creds, err := s.repo.GetAuthByUsername(ctx, name)
	if err != nil {
		return errs.ErrUserNotFound
	}

	if creds.PasswordHash != password {
		return errors.New("invalid password")
	}

	return nil
}

func (s *UserService) Create(ctx context.Context, user *models.User, password string) error {
	if err := user.Validate(); err != nil {
		return err
	}
	if password == "" {
		return errors.New("password cannot be empty")
	}

	
	err := s.repo.Create(ctx, *user, password)
	if err != nil {
		return fmt.Errorf("s.repo.Create: %w", err)
	}

	return nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, errs.ErrValidation
	}
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	err := s.repo.Update(ctx, *user)
	if err != nil {
		return fmt.Errorf("s.repo.Update: %w", err)
	}

	return nil
}

func (s *UserService) SoftDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return errs.ErrValidation
	}

	return s.repo.SoftDelete(ctx, id)
}
