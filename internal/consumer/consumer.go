package consumer

import (
	"context"
	events "project/internal/service/eventBus"
	"project/pkg/logger"
	"sync"
	"go.uber.org/zap"
)

func StartAuditConsumer(ctx context.Context, wg *sync.WaitGroup, bus *events.Bus, log *logger.Logger){
	wg.Add(1)
	go func ()  {
			defer wg.Done()
			for {
				select{
				case <-ctx.Done():
					return
				case event:=<-bus.Subscribe():
					log.Audit.Info(
						event.Type, 
						zap.Int("user_id", event.UserID),
					)
				}
			}
	}()
}