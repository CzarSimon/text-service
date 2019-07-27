package main

import (
	"database/sql"
	"fmt"

	"github.com/CzarSimon/text-service/go/pkg/repository"
	"github.com/CzarSimon/text-service/go/pkg/service"
	"github.com/CzarSimon/text-service/go/pkg/utils/dbutil"
	"github.com/CzarSimon/text-service/go/pkg/utils/environ"
	"go.uber.org/zap"
)

type env struct {
	cfg        config
	db         *sql.DB
	textGetter service.TextGetter
}

func (e *env) Close() error {
	return e.db.Close()
}

func getEnv(cfg config) *env {
	db := dbutil.MustConnect(cfg.db)

	err := dbutil.Upgrade(cfg.migrationsPath, cfg.db.Driver(), db)
	if err != nil {
		log.Panic("Failed to apply migratons", zap.Error(err))
	}

	languageRepo := repository.NewLanguageRepository(db)
	textRepo := repository.NewTextRepository(db)

	return &env{
		cfg:        cfg,
		db:         db,
		textGetter: service.NewTextGetter(languageRepo, textRepo),
	}
}

func (e *env) checkHealth() error {
	return dbutil.Connected(e.db)
}

type config struct {
	db             dbutil.Config
	port           string
	migrationsPath string
}

func getConfig() config {
	storageType := environ.Get("STORAGE", "posgres")
	baseMigrations := environ.Get("MIGRATIONS_PATH", "/etc/text-service/migrations")

	return config{
		db:             getDBConfig(storageType),
		port:           environ.Get("SERVICE_PORT", "8080"),
		migrationsPath: fmt.Sprintf("%s/%s", baseMigrations, storageType),
	}
}

func getDBConfig(storageType string) dbutil.Config {
	switch storageType {
	case "sqlite":
		return dbutil.SqliteConfig{Name: environ.MustGet("DB_NAME")}
	case "memory":
		return dbutil.SqliteConfig{}
	default:
		return dbutil.SqliteConfig{}
	}
}
