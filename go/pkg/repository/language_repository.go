package repository

import (
	"database/sql"
	"time"

	"github.com/CzarSimon/text-service/go/pkg/models"
	"github.com/CzarSimon/text-service/go/pkg/utils/context"
	"github.com/CzarSimon/text-service/go/pkg/utils/logger"
	"github.com/pkg/errors"
)

var log = logger.GetDefaultLogger("pkg/repository").Sugar()

// Common errors
var (
	ErrNotFound = errors.New("not found")
)

// LanguageRepository storage interface for languages.
type LanguageRepository interface {
	Find(ctx *context.Context, language string) (models.Language, error)
	Save(ctx *context.Context, language string) error
}

// NewLanguageRepository creates a new LanguageRepository using the default implementation.
func NewLanguageRepository(db *sql.DB) LanguageRepository {
	return &languageRepo{
		db: db,
	}
}

type languageRepo struct {
	db *sql.DB
}

const findLanguagesQuery = `SELECT id, created_at FROM language WHERE id = $1`

func (r *languageRepo) Find(ctx *context.Context, language string) (models.Language, error) {
	log.Debugw("languageRepo.Find", "language", language, "ctx", ctx)

	var lang models.Language
	err := r.db.QueryRowContext(ctx, findLanguagesQuery, language).Scan(&lang.ID, &lang.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Language{}, ErrNotFound
	}

	if err != nil {
		return models.Language{}, errors.Wrapf(err, "Failed to query language. language=%s", language)
	}

	return lang, nil
}

const saveLanguageQuery = `INSERT INTO language(id, created_at) VALUES ($1, $2)`

func (r *languageRepo) Save(ctx *context.Context, language string) error {
	log.Debugw("languageRepo.Save", "language", language, "ctx", ctx)
	_, err := r.db.ExecContext(ctx, saveLanguageQuery, language, time.Now())
	if err != nil {
		return errors.Wrapf(err, "Failed to insert language. language=%s", language)
	}

	return nil
}
