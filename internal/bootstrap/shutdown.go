package bootstrap

import (
	"context"
	"fmt"
	"log"
)

func (app *App) Shutdown(ctx context.Context) error {
	log.Println("🛑 Shutting down server...")

	if err := app.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("🛑 Closing database connection...")
	if err := app.DB.Close(); err != nil {
		return fmt.Errorf("database close failed: %w", err)
	}

	log.Println("✅ App shutdown complete")
	return nil
}
