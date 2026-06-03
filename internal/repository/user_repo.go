package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"project/internal/models"
	"project/pkg/errs"
	"sync"
)


type UserRepo struct{
	mu sync.Mutex
	fileName string
}

func New(fileName string) *UserRepo{
	return &UserRepo{
		mu: sync.Mutex{},
		fileName: fileName,
	}
}

func(s *UserRepo) getUsers(ctx context.Context)([]models.User, error){
	file, err:=os.ReadFile(s.fileName)
	if err!=nil{
		return nil, fmt.Errorf("os.Readfile:%w", err)
	}
	var users []models.User
	err=json.Unmarshal(file, &users)
	if err!=nil{
		return nil, fmt.Errorf("json.Unmarshal:%w", err)
	}
	return users, nil
}

func (s *UserRepo)GetAll(ctx context.Context)([]models.User, error){
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getUsers(ctx)
}

func (s *UserRepo)GetById(ctx context.Context, id int)(*models.User, error){
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err:=s.getUsers(ctx)
	if err!=nil{
		return nil, fmt.Errorf("s.getUsers:%w", err)
	}
	var user models.User
	for _, value:=range users{
		if value.ID==id{
			user=value
			return &user, nil
		}
	}
	return nil, errs.ErrUserNotFound
}

func(s *UserRepo)Create(ctx context.Context, user models.User)error{
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err:=s.getUsers(ctx)
	if err!=nil{
		return fmt.Errorf("s.getUsers:%w", err)
	}
	users=append(users, user)
	file, err:=json.Marshal(users)
	if err!=nil{
		return fmt.Errorf("json.Marshal:%w", err)
	}
	err=os.WriteFile(s.fileName, file, 0644)
	if err!=nil{
		return fmt.Errorf("os.Writefile:%w", err)
	}
	return nil
}

func(s *UserRepo)Update(ctx context.Context, user models.User)error{
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err:=s.getUsers(ctx)
	if err!=nil{
		return fmt.Errorf("s.getUsers:%w", err)
	}
	for i, value:=range users{
		if value.ID==user.ID{
			users[i]=user
		}
	}
	fileData, err:=json.Marshal(users)
	if err!=nil{
		return fmt.Errorf("json.Marshal:%w", err)
	}
	err=os.WriteFile(s.fileName, fileData, 0644)
	if err!=nil{
		return fmt.Errorf("os.WriteFile:%w", err)
	}
	return nil
}

