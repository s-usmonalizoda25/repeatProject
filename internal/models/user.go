package models

import (
	"fmt"
	"project/pkg/errs"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}

func (u *User) Validate() error {
	if u.ID <= 0 {
		return fmt.Errorf("user id: %w", errs.ErrValidation)
	}
	if u.Age <= 0 {
		return fmt.Errorf("user age: %w", errs.ErrValidation)
	}
	if u.Name == "" {
		return fmt.Errorf("user name: %w", errs.ErrValidation)
	}
	if u.Password == "" {
		return fmt.Errorf("user password: %w", errs.ErrValidation)
	}
	return nil
}
