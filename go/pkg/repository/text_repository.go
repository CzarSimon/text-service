package repository

import (
	"database/sql"
	"time"

	"github.com/CzarSimon/text-service/go/pkg/models"
	"github.com/CzarSimon/text-service/go/pkg/utils/context"
	"github.com/pkg/errors"
)

// TextRepository storage interface for translated texts.
type TextRepository interface {
	Find(ctx *context.Context, key, language string) (models.TranslatedText, error)
	Save(ctx *context.Context, text models.TranslatedText) error
}

// NewTextRepository creates a new TextRepository using the default implementation.
func NewTextRepository(db *sql.DB) TextRepository {
	return &textRepo{
		db: db,
	}
}

type textRepo struct {
	db *sql.DB
}

const findTextQuery = `SELECT id, key, language, value, created_at, updated_at FROM translated_text WHERE key = $1 AND language = $2`

func (r *textRepo) Find(ctx *context.Context, key, language string) (models.TranslatedText, error) {
	log.Debugw("textRepo.Find", "key", key, "language", language, "ctx", ctx)

	var t models.TranslatedText
	err := r.db.QueryRowContext(ctx, findTextQuery, key, language).Scan(&t.ID, &t.Key, &t.Language, &t.Value, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return models.TranslatedText{}, ErrNotFound
	}

	if err != nil {
		return models.TranslatedText{}, errors.Wrapf(err, "Failed to query translated_text: %s", language)
	}

	return t, nil
}

const saveTextQuery = `INSERT INTO translated_text(key, language, value, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`

func (r *textRepo) Save(ctx *context.Context, text models.TranslatedText) error {
	log.Debugw("textRepo.Save", "key", text.Key, "language", text.Language, "ctx", ctx)

	now := time.Now()
	_, err := r.db.ExecContext(ctx, saveTextQuery, text.Key, text.Language, text.Value, now, now)
	if err != nil {
		return errors.Wrapf(err, "Failed to insert translated_text: key=%s", text.Key)
	}

	return nil
}
