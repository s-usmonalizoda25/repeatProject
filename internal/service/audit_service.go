package service

import (
	"context"
	"project/internal/models"
	"project/internal/repository"
)

type AuditService struct {
	repo *repository.AuditRepo
}

func NewAuditService(repo *repository.AuditRepo) *AuditService {
	return &AuditService{
		repo: repo,
	}
}

func (s *AuditService) Record(ctx context.Context, action string, userID int) error {
	entry := models.AuditEntry{
		Action: action,
		UserID: userID,
	}
	return s.repo.Save(ctx, entry)
}
