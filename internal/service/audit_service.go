package service

import (
	"context"
	"project/internal/models"
	"project/internal/repository"
	"time"
)

type AuditService struct{
	repo *repository.AuditRepo
}

func NewAuditService(repo *repository.AuditRepo)*AuditService{
	return &AuditService{
		repo: repo,
	}
}
func (s *AuditService) Record(ctx context.Context, action string, user_id int, name string, age int)error{
	entry:=models.Audit{
		Action: action,
		UserId: user_id,
		Time: time.Now(),
		UserName: name,
		UserAge: age,
	}
	return s.repo.Log(ctx, entry)
}