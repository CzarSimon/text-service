package main

import (
	stdctx "context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CzarSimon/text-service/go/pkg/models"
	"github.com/CzarSimon/text-service/go/pkg/repository"
	"github.com/CzarSimon/text-service/go/pkg/utils/context"
	"github.com/CzarSimon/text-service/go/pkg/utils/dbutil"
	"github.com/CzarSimon/text-service/go/pkg/utils/httputil"
	"github.com/CzarSimon/text-service/go/pkg/utils/id"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

/*
   en_headers = headers(language="en")
   res_missing = client.get("/v1/texts/key/ONLY_SV_TEXT_KEY", headers=en_headers)
   assert res_missing.status_code == status.HTTP_404_NOT_FOUND
*/

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

func TestGetTextNoAndUnsupportedLanguage(t *testing.T) {
	assert := assert.New(t)
	e := createTestEnv()
	server := newServer(e)

	// Swedish text value test
	path := "/v1/texts/key/TEST_TEXT_KEY"
	req := createTestRequest(path, "")
	res := performTestRequest(server.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)

	req = createTestRequest(path, "xy")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)
}

/*
def test_get_text_by_group():
    items = [
        Language(id="sv"),
        Language(id="en"),
        TranslatedText(key="TEST_TEXT_KEY", language="sv", value="sv-text-val"),
        TranslatedText(key="TEST_TEXT_KEY", language="en", value="en-text-val"),
        TranslatedText(key="OTHER_TEXT_KEY", language="sv", value="sv-other-val"),
        TranslatedText(key="OTHER_TEXT_KEY", language="en", value="en-other-val"),
        TranslatedText(key="NOT_IN_GROUP", language="sv", value="sv-other-val"),
        TranslatedText(key="NOT_IN_GROUP", language="en", value="en-other-val"),
        TranslatedText(key="ONLY_SV_TEXT_KEY", language="sv", value="sv-only-val"),
        Group(id="MOBILE_APP"),
        TextGroup(text_key="TEST_TEXT_KEY", group_id="MOBILE_APP"),
        TextGroup(text_key="OTHER_TEXT_KEY", group_id="MOBILE_APP"),
        TextGroup(text_key="ONLY_SV_TEXT_KEY", group_id="MOBILE_APP"),
    ]

    with TestEnvironment(items) as client:
        sv_headers_ok = headers(language="sv")
        res_sv = client.get("/v1/texts/group/MOBILE_APP", headers=sv_headers_ok)
        assert res_sv.status_code == status.HTTP_200_OK
        body = res_sv.get_json()
        assert len(body) == 3
        assert body["TEST_TEXT_KEY"] == "sv-text-val"
        assert body["OTHER_TEXT_KEY"] == "sv-other-val"
        assert body["ONLY_SV_TEXT_KEY"] == "sv-only-val"

        no_lang_headers = headers(language=None)
        res_no_lang = client.get("/v1/texts/group/MOBILE_APP", headers=no_lang_headers)
        assert res_no_lang.status_code == status.HTTP_400_BAD_REQUEST

        en_headers = headers(language="en")
        res_missing = client.get("/v1/texts/group/MISSING_GROUP", headers=en_headers)
        assert res_missing.status_code == status.HTTP_404_NOT_FOUND

        wrong_lang_headers = headers(language="xy")
        res_wrong_lang = client.get(
            "/v1/texts/group/MOBILE_APP", headers=wrong_lang_headers
        )
        assert res_wrong_lang.status_code == status.HTTP_400_BAD_REQUEST

        en_headers_ok = headers(language="en")
        res_en = client.get("/v1/texts/group/MOBILE_APP", headers=en_headers_ok)
        assert res_en.status_code == status.HTTP_200_OK
        body = res_en.get_json()
        assert len(body) == 3
        assert body["TEST_TEXT_KEY"] == "en-text-val"
        assert body["OTHER_TEXT_KEY"] == "en-other-val"
		assert body["ONLY_SV_TEXT_KEY"] == None
*/

/*
func TestGetTextByGroup(t *testing.T) {
	assert := assert.New(t)
	e := createTestEnv()
	server := newServer(e)

	path := "/v1/texts/group/"
	req := createTestRequest(path, "sv")
	res := performTestRequest(server.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)
}
*/

func TestGetTextByGroupNoAndUnsupportedLanguage(t *testing.T) {
	assert := assert.New(t)
	e := createTestEnv()
	server := newServer(e)

	// Swedish text value test
	path := "/v1/texts/group/MOBILE_APP"
	req := createTestRequest(path, "")
	res := performTestRequest(server.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)
}

func createTestEnv() *env {
	cfg := config{
		db:             dbutil.SqliteConfig{},
		migrationsPath: "../resources/db/sqlite",
	}

	e := getEnv(cfg)
	langRepo := repository.NewLanguageRepository(e.db)
	textRepo := repository.NewTextRepository(e.db)

	ctx := context.New(stdctx.Background(), "createTestEnv", "")
	errs := make([]error, 7)
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
		Key: "ONLY_SV_TEXT_KEY", Language: "sv", Value: "en-only-val",
	})
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
