package consumer

import (
	"context"
	"sync"
	"time"

	"project/internal/service"
	"project/internal/service/eventBus"
	"project/pkg/logger"

	"go.uber.org/zap"
)

func StartAuditConsumer(ctx context.Context, wg *sync.WaitGroup, bus *eventBus.Bus, loggy *logger.Logger, auditService *service.AuditService) {
	wg.Add(1)

	ch := bus.Subscribe()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				loggy.Info("Audit consumer stopping...")
				return
			case event, ok := <-ch:
				if !ok {
					return
				}

				loggy.AuditLogger().Info("Audit Event",
					zap.String("action", event.Type),
					zap.Int("user_id", event.UserID),
					zap.Time("timestamp", time.Now()),
				)

				entryCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				_ = auditService.Record(entryCtx, event.Type, event.UserID)
				cancel()
			}
		}
	}()
}
