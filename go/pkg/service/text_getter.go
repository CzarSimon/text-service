package service

import (
	"fmt"

	"github.com/CzarSimon/text-service/go/pkg/models"
	"github.com/CzarSimon/text-service/go/pkg/repository"
	"github.com/CzarSimon/text-service/go/pkg/utils/context"
	"github.com/CzarSimon/text-service/go/pkg/utils/httputil"
	"github.com/CzarSimon/text-service/go/pkg/utils/logger"
)

var log = logger.GetDefaultLogger("pkg/service").Sugar()

// TextGetter interface for getting texts by for a given language.
type TextGetter interface {
	Get(ctx *context.Context, key string) (models.Texts, error)
	GetGroup(ctx *context.Context, groupID string) (models.Texts, error)
}

// NewTextGetter creates a new TextGetter using the default implementation.
func NewTextGetter(
	languageRepo repository.LanguageRepository,
	textRepo repository.TextRepository,
	groupRepo repository.GroupRepository) TextGetter {
	return &getter{
		languageRepo: languageRepo,
		textRepo:     textRepo,
		groupRepo:    groupRepo,
	}
}

type getter struct {
	languageRepo repository.LanguageRepository
	textRepo     repository.TextRepository
	groupRepo    repository.GroupRepository
}

func (g *getter) Get(ctx *context.Context, key string) (models.Texts, error) {
	log.Debugw("getter.Get", "key", key, "ctx", ctx)
	err := g.assertLanguageExists(ctx)
	if err != nil {
		return nil, err
	}

	text, err := g.textRepo.Find(ctx, key, ctx.Language)
	if err == repository.ErrNotFound {
		return nil, httputil.ErrNotFound
	}
	if err != nil {
		log.Errorw("Failed to find text by key", "error", err, "ctx", ctx)
		return nil, httputil.ErrInternalServerError
	}

	return mapTextsToMap(text), nil
}

func (g *getter) GetGroup(ctx *context.Context, groupID string) (models.Texts, error) {
	log.Debugw("getter.GetGroup", "groupId", groupID, "ctx", ctx)

	err := g.assertLanguageExists(ctx)
	if err != nil {
		return nil, err
	}

	texts, err := g.groupRepo.FindTexts(ctx, groupID, ctx.Language)
	if err != nil {
		log.Errorw("Failed to find texts by group", "error", err, "ctx", ctx)
		return nil, httputil.ErrInternalServerError
	}

	if len(texts) == 0 {
		log.Infow("No texts for group. groupId="+groupID, "ctx", ctx)
		return nil, httputil.ErrNotFound
	}

	return mapTextsToMap(texts...), nil
}

func (g *getter) assertLanguageExists(ctx *context.Context) error {
	_, err := g.languageRepo.Find(ctx, ctx.Language)
	if err == repository.ErrNotFound {
		log.Infow("Language not found", "ctx", ctx)
		errorMsg := fmt.Sprintf("Unsupported language: %s", ctx.Language)
		return httputil.BadRequest(errorMsg)
	}

	if err != nil {
		log.Errorw("Failed to find if language is supported", "error", err, "ctx", ctx)
		return httputil.ErrInternalServerError
	}

	return nil
}

func mapTextsToMap(texts ...models.TranslatedText) models.Texts {
	textMap := make(models.Texts)
	for _, text := range texts {
		textMap[text.Key] = text.Value
	}

	return textMap
}
