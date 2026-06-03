package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"project/internal/models"
	"sync"
)

type AuditRepo struct{
	mu sync.Mutex
	fileName string
}

func NewAuditRepo(fileName string)*AuditRepo{
	return &AuditRepo{
		mu: sync.Mutex{},
		fileName: fileName,
	}
}


func (r *AuditRepo) Log(ctx context.Context, entry models.Audit)error{
	r.mu.Lock()
	defer r.mu.Unlock()
	var entries []models.Audit

	file, err:=os.ReadFile(r.fileName)
	if err==nil{
		json.Unmarshal(file, &entries)
	}
	entries=append(entries, entry)
	fileData, err:=json.MarshalIndent(entries, " ", "   ")
	if err!=nil{
		return fmt.Errorf("json.Unmarshal:%w", err)
	}
	return os.WriteFile(r.fileName, fileData, 0644)
}