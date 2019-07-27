package main

import (
	stdctx "context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CzarSimon/text-service/go/pkg/models"
	"github.com/CzarSimon/text-service/go/pkg/repository"
	"github.com/CzarSimon/text-service/go/pkg/utils/context"
	"github.com/CzarSimon/text-service/go/pkg/utils/httputil"
	"github.com/CzarSimon/text-service/go/pkg/utils/id"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetTextOK(t *testing.T) {
	assert := assert.New(t)
	e := createTestEnv()
	server := newServer(e)

	// Swedish text value test
	path := "/v1/texts/key/TEST_TEXT_KEY"
	req := createTestRequest(path, "sv")
	res := performTestRequest(server.Handler, req)
	assert.Equal(http.StatusOK, res.Code)

	var svtexts models.Texts
	err := json.NewDecoder(res.Body).Decode(&svtexts)
	assert.NoError(err)

	assert.Equal(1, len(svtexts))
	text, ok := svtexts["TEST_TEXT_KEY"]
	assert.True(ok)
	assert.Equal("sv-text-val", text)

	// English text value test
	req = createTestRequest(path, "en")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusOK, res.Code)

	var entexts models.Texts
	err = json.NewDecoder(res.Body).Decode(&entexts)
	assert.NoError(err)

	assert.Equal(1, len(entexts))
	text, ok = entexts["TEST_TEXT_KEY"]
	assert.True(ok)
	assert.Equal("en-text-val", text)
}

func TestGetTextFail(t *testing.T) {
	assert := assert.New(t)
	e := createTestEnv()
	server := newServer(e)

	// No language
	path := "/v1/texts/key/TEST_TEXT_KEY"
	req := createTestRequest(path, "")
	res := performTestRequest(server.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)

	// Unsupported language
	req = createTestRequest(path, "xy")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)

	// No such text for the given language
	path = "/v1/texts/key/ONLY_SV_TEXT_KEY"
	req = createTestRequest(path, "en")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusNotFound, res.Code)
}

/*
   	en_headers = headers(language="en")
	res_missing = client.get("/v1/texts/group/MISSING_GROUP", headers=en_headers)
    assert res_missing.status_code == status.HTTP_404_NOT_FOUND
*/

func TestGetTextsByGroupOK(t *testing.T) {
	assert := assert.New(t)
	e := createTestEnv()
	server := newServer(e)

	// Swedish texts in group test.
	path := "/v1/texts/group/MOBILE_APP"
	req := createTestRequest(path, "sv")
	res := performTestRequest(server.Handler, req)
	assert.Equal(http.StatusOK, res.Code)

	var svtexts models.Texts
	err := json.NewDecoder(res.Body).Decode(&svtexts)
	assert.NoError(err)

	assert.Equal(3, len(svtexts))
	text, _ := svtexts["TEST_TEXT_KEY"]
	assert.Equal("sv-text-val", text)
	text, _ = svtexts["OTHER_TEXT_KEY"]
	assert.Equal("sv-other-val", text)
	text, _ = svtexts["ONLY_SV_TEXT_KEY"]
	assert.Equal("sv-only-val", text)

	// English texts in group test.
	req = createTestRequest(path, "en")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusOK, res.Code)

	var entexts models.Texts
	err = json.NewDecoder(res.Body).Decode(&entexts)
	assert.NoError(err)

	assert.Equal(2, len(entexts))
	text, _ = entexts["TEST_TEXT_KEY"]
	assert.Equal("en-text-val", text)
	text, _ = entexts["OTHER_TEXT_KEY"]
	assert.Equal("en-other-val", text)
	_, ok := entexts["ONLY_SV_TEXT_KEY"]
	assert.False(ok)
}

func TestGetTextsByGroupFail(t *testing.T) {
	assert := assert.New(t)
	e := createTestEnv()
	server := newServer(e)

	// No language
	path := "/v1/texts/group/MOBILE_APP"
	req := createTestRequest(path, "")
	res := performTestRequest(server.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)

	// Unsupported language
	req = createTestRequest(path, "xy")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)

	// Missing group swedish
	path = "/v1/texts/group/MISSING_GROUP"
	req = createTestRequest(path, "sv")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusNotFound, res.Code)

	// Missing group english
	path = "/v1/texts/group/MISSING_GROUP"
	req = createTestRequest(path, "en")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusNotFound, res.Code)
}

func createTestEnv() *env {
	os.Setenv("STORAGE", "memory")
	os.Setenv("MIGRATIONS_PATH", "../resources/db")

	cfg := getConfig()
	e := getEnv(cfg)

	langRepo := repository.NewLanguageRepository(e.db)
	textRepo := repository.NewTextRepository(e.db)
	groupRepo := repository.NewGroupRepository(e.db)

	ctx := context.New(stdctx.Background(), "createTestEnv", "")
	errs := make([]error, 13)
	errs[0] = langRepo.Save(ctx, "sv")
	errs[1] = langRepo.Save(ctx, "en")
	errs[2] = textRepo.Save(ctx, models.TranslatedText{
		Key: "TEST_TEXT_KEY", Language: "sv", Value: "sv-text-val",
	})
	errs[3] = textRepo.Save(ctx, models.TranslatedText{
		Key: "TEST_TEXT_KEY", Language: "en", Value: "en-text-val",
	})
	errs[4] = textRepo.Save(ctx, models.TranslatedText{
		Key: "OTHER_TEXT_KEY", Language: "sv", Value: "sv-other-val",
	})
	errs[5] = textRepo.Save(ctx, models.TranslatedText{
		Key: "OTHER_TEXT_KEY", Language: "en", Value: "en-other-val",
	})
	errs[6] = textRepo.Save(ctx, models.TranslatedText{
		Key: "ONLY_SV_TEXT_KEY", Language: "sv", Value: "sv-only-val",
	})
	errs[7] = textRepo.Save(ctx, models.TranslatedText{
		Key: "NOT_IN_GROUP", Language: "sv", Value: "sv-non-group-val",
	})
	errs[8] = textRepo.Save(ctx, models.TranslatedText{
		Key: "NOT_IN_GROUP", Language: "en", Value: "en-non-group-val",
	})
	errs[9] = groupRepo.Save(ctx, models.TextGroup{ID: "MOBILE_APP"})
	errs[10] = groupRepo.AddTextToGroup(ctx, "TEST_TEXT_KEY", "MOBILE_APP")
	errs[11] = groupRepo.AddTextToGroup(ctx, "OTHER_TEXT_KEY", "MOBILE_APP")
	errs[12] = groupRepo.AddTextToGroup(ctx, "ONLY_SV_TEXT_KEY", "MOBILE_APP")
	ensureNoErrors(errs)

	return e
}

func performTestRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func createTestRequest(route, language string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		log.Fatal("Failed to create request", zap.Error(err))
	}

	req.Header.Set(httputil.RequestIDHeader, id.New())
	req.Header.Set(httputil.AcceptLanguage, language)
	return req
}

func ensureNoErrors(errs []error) {
	for i, err := range errs {
		if err != nil {
			log.Fatal(fmt.Sprintf("%d - Error: %s", i, err))
		}
	}
}
