package service

import (
	"context"
	"project/internal/models"
	"project/internal/repository"
	"project/pkg/errs"
)

type UserService struct{
	repo *repository.UserRepo
	AuditService *AuditService
}

func New(repo *repository.UserRepo, auditservice *AuditService)*UserService{
	return &UserService{
		repo: repo,
		AuditService: auditservice,
	}
}
func (s *UserService) GetAll(ctx context.Context)([]models.User, error){
	return s.repo.GetAll(ctx)
}

func(s *UserService) Create(ctx context.Context, user *models.User)error{
	if err:=user.Validate();err!=nil{
		return errs.ErrValidation
	}

	err:=s.repo.Create(ctx, *user)
	if err!=nil{
		return err
	}
	s.AuditService.Record(ctx, "user_created", user.ID, user.Name, user.Age)
	return nil
}

func(s *UserService) Update(ctx context.Context, user *models.User)error{
	if err:=user.Validate(); err!=nil{
		return errs.ErrValidation
	}
	err:=s.repo.Update(ctx, *user)
	if err!=nil{
		return err
	}
	s.AuditService.Record(ctx, "user_created", user.ID, user.Name, user.Age)
	return nil
}

func (s *UserService)GetById(ctx context.Context, id int)(*models.User, error){
	if id<=0{
		return nil, errs.ErrValidation
	}
	return s.repo.GetById(ctx, id)
}


