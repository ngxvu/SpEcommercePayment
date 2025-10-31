package bootstrap

import (
	"fmt"
	repo "payment/internal/repositories/pg-gorm"
	"payment/pkg/core/configloader"
	"payment/pkg/core/db"
)

type App struct {
	Config *configloader.Config
	PGRepo repo.PGInterface
}

// initializeApp initializes all application dependencies
func InitializeApp() (*App, error) {
	config := configloader.GetConfig()

	// Initialize database
	dbBackend, err := db.InitDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	pgRepo := repo.NewPGRepo(dbBackend)

	return &App{
		PGRepo: pgRepo,
		Config: config,
	}, nil
}
