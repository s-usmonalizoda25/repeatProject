package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"project/internal/models"
	"sync"
)

type AuditRepo struct {
	mu       sync.Mutex
	fileName string
}

func NewAuditRepo(fileName string) *AuditRepo {
	return &AuditRepo{
		fileName: fileName,
	}
}

func (r *AuditRepo) Save(ctx context.Context, entry models.AuditEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	var entries []models.AuditEntry
	file, err := os.ReadFile(r.fileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	if len(file) > 0 {
		if err := json.Unmarshal(file, &entries); err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	entries = append(entries, entry)

	fileData, err := json.MarshalIndent(entries, "", "    ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	if err := os.WriteFile(r.fileName, fileData, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}

	return nil
}
