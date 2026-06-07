package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"project/internal/models"
	"project/pkg/errs"
	"sync"
)

type UserRepo struct {
	mu       sync.Mutex
	fileName string
}

func New(fileName string) *UserRepo {
	return &UserRepo{
		fileName: fileName,
	}
}

func (s *UserRepo) getUsers(ctx context.Context) ([]models.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	file, err := os.ReadFile(s.fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.User{}, nil
		}
		return nil, fmt.Errorf("os.Readfile:%w", err)
	}
	if len(file) == 0 {
		return []models.User{}, nil
	}
	var users []models.User
	err = json.Unmarshal(file, &users)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal:%w", err)
	}
	return users, nil
}

func (s *UserRepo) saveUsers(ctx context.Context, users []models.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	fileData, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent:%w", err)
	}
	err = os.WriteFile(s.fileName, fileData, 0644)
	if err != nil {
		return fmt.Errorf("os.WriteFile:%w", err)
	}
	return nil
}

func (s *UserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getUsers(ctx)
}

func (s *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err := s.getUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, value := range users {
		if value.ID == id {
			return &value, nil
		}
	}
	return nil, errs.ErrUserNotFound
}

func (s *UserRepo) GetByName(ctx context.Context, name string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err := s.getUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, value := range users {
		if value.Name == name {
			return &value, nil
		}
	}
	return nil, errs.ErrUserNotFound
}

func (s *UserRepo) Create(ctx context.Context, user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err := s.getUsers(ctx)
	if err != nil {
		return err
	}

	for _, value := range users {
		if value.ID == user.ID {
			return fmt.Errorf("user with id %d already exists: %w", user.ID, errs.ErrValidation)
		}
	}

	user.IsActive = true
	users = append(users, user)
	return s.saveUsers(ctx, users)
}

func (s *UserRepo) Update(ctx context.Context, user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err := s.getUsers(ctx)
	if err != nil {
		return err
	}
	found := false
	for i, value := range users {
		if value.ID == user.ID {
			users[i] = user
			found = true
			break
		}
	}
	if !found {
		return errs.ErrUserNotFound
	}
	return s.saveUsers(ctx, users)
}

func (s *UserRepo) SoftDelete(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err := s.getUsers(ctx)
	if err != nil {
		return err
	}

	for i, value := range users {
		if value.ID == id {
			if !value.IsActive {
				return errors.New("user is already deactivated")
			}
			users[i].IsActive = false
			return s.saveUsers(ctx, users)
		}
	}
	return errs.ErrUserNotFound
}
