package service

import (
	"context"
	"fmt"

	"project/internal/models"
	"project/internal/repository"
	"project/internal/service/eventBus"
	"project/pkg/errs"
)

type UserService struct {
	repo *repository.UserRepo
	bus  *eventBus.Bus
}

func New(repo *repository.UserRepo, bus *eventBus.Bus) *UserService {
	return &UserService{
		repo: repo,
		bus:  bus,
	}
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetAll: %w", err)
	}

	s.bus.Publish(eventBus.Event{
		Type:   "Get all users",
		UserID: 10,
	})

	return users, nil
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return errs.ErrValidation
	}

	err := s.repo.Create(ctx, *user)
	if err != nil {
		return fmt.Errorf("s.repo.Create: %w", err)
	}

	s.bus.Publish(eventBus.Event{
		Type:   "Created User",
		UserID: user.ID,
	})

	return nil
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return errs.ErrValidation
	}
	
	err := s.repo.Update(ctx, *user)
	if err != nil {
		return fmt.Errorf("s.repo.Update: %w", err)
	}

	s.bus.Publish(eventBus.Event{
		Type:   "Updated User",
		UserID: user.ID,
	})

	return nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, errs.ErrValidation
	}
	return s.repo.GetById(ctx, id)
}