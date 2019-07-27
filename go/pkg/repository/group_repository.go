package repository

import (
	"database/sql"
	"time"

	"github.com/CzarSimon/text-service/go/pkg/models"
	"github.com/CzarSimon/text-service/go/pkg/utils/context"
	"github.com/pkg/errors"
)

// GroupRepository storage interface for translated texts.
type GroupRepository interface {
	FindTexts(ctx *context.Context, groupID, language string) ([]models.TranslatedText, error)
	Save(ctx *context.Context, group models.TextGroup) error
	AddTextToGroup(ctx *context.Context, textKey, groupID string) error
}

// NewGroupRepository creates a new GroupRepository using the default implementation.
func NewGroupRepository(db *sql.DB) GroupRepository {
	return &groupRepo{
		db: db,
	}
}

type groupRepo struct {
	db *sql.DB
}

const findGroupTextsQuery = `
	SELECT t.id, t.key, t.language, t.value, t.created_at, t.updated_at 
	FROM translated_text t
	INNER JOIN text_group_membership tgm ON t.key = tgm.text_key
	WHERE tgm.group_id = $1 AND t.language = $2`

func (r *groupRepo) FindTexts(ctx *context.Context, groupID, language string) ([]models.TranslatedText, error) {
	log.Debugw("groupRepo.FindTexts", "groupId", groupID, "language", language, "ctx", ctx)
	rows, err := r.db.QueryContext(ctx, findGroupTextsQuery, groupID, language)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to query group texts. groupId=%s language=%s", groupID, language)
	}
	defer rows.Close()

	texts := make([]models.TranslatedText, 0)
	var t models.TranslatedText
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Key, &t.Language, &t.Value, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to scan text. groupId=%s language=%s", groupID, language)
		}

		texts = append(texts, t)
	}

	return texts, nil
}

const saveGroupQuery = `INSERT INTO text_group(id, created_at) VALUES ($1, $2)`

func (r *groupRepo) Save(ctx *context.Context, group models.TextGroup) error {
	log.Debugw("groupRepo.Save", "groupId", group.ID, "ctx", ctx)
	_, err := r.db.ExecContext(ctx, saveGroupQuery, group.ID, time.Now())
	if err != nil {
		return errors.Wrapf(err, "Failed to save group. id=%s", group.ID)
	}

	return nil
}

const addTextToGroupQuery = `INSERT INTO text_group_membership(text_key, group_id, created_at) VALUES ($1, $2, $3)`

func (r *groupRepo) AddTextToGroup(ctx *context.Context, textKey, groupID string) error {
	log.Debugw("groupRepo.AddTextToGroup", "textKey", textKey, "groupId", groupID, "ctx", ctx)
	_, err := r.db.ExecContext(ctx, addTextToGroupQuery, textKey, groupID, time.Now())
	if err != nil {
		return errors.Wrapf(err, "Failed to add text to group. textKey=%s groupId=%s", textKey, groupID)
	}

	return nil
}
