package main

import (
	"context"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/logging"
	"time"
)

const (
	userGameTTL      = 30 * 24 * time.Hour
	anonymousGameTTL = 24 * time.Hour
)

func startPruneWorker(ctx context.Context, interval time.Duration) {
	go func() {
		prune := func() {
			if err := db.PruneGameHistory(ctx, userGameTTL, anonymousGameTTL); err != nil {
				logging.Error("Failed to prune game history", err)
			}
		}
		prune()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				prune()
			}
		}
	}()
}
