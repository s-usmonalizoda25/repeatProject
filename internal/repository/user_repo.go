package repository

import (
	"context"
	"database/sql"
	"errors"
	"project/internal/models"
	"project/pkg/errs"
)

type IUserRepo interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetAuthByUsername(ctx context.Context, username string) (*models.Credentials, error)
	Create(ctx context.Context, user models.User, password string) error
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
	const query = `SELECT id, name, age, is_active FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	const query = `SELECT id, name, age, is_active FROM users WHERE id = $1`

	var u models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Age, &u.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetAuthByUsername(ctx context.Context, username string) (*models.Credentials, error) {
	const query = `SELECT id, user_id, username, password_hash FROM auth WHERE username = $1`

	var creds models.Credentials
	err := r.db.QueryRowContext(ctx, query, username).Scan(&creds.ID, &creds.UserID, &creds.Username, &creds.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return &creds, nil
}

func (r *userRepository) Create(ctx context.Context, user models.User, password string) error {
	const (
		userQuery = `INSERT INTO users (name, age, is_active) VALUES ($1, $2, $3) RETURNING id`
		authQuery = `INSERT INTO auth (user_id, username, password_hash) VALUES ($1, $2, $3)`
	)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRowContext(ctx, userQuery, user.Name, user.Age, user.IsActive).Scan(&userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, authQuery, userID, user.Name, password)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *userRepository) Update(ctx context.Context, user models.User) error {
	const query = `UPDATE users SET name = $1, age = $2, is_active = $3 WHERE id = $4`

	res, err := r.db.ExecContext(ctx, query, user.Name, user.Age, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) SoftDelete(ctx context.Context, id int) error {
	const query = `UPDATE users SET is_active = false WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}
