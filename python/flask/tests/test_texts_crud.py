# Standard library
import json
from datetime import datetime

# 3rd party modules
from app.httputil import status

# Intenal modules
from tests import TestEnvironment, JSON, headers, new_id
from app.models import TranslatedText, Language, TextGroup, TextGroupMembership


def test_get_text():
    items = [
        Language(id="sv"),
        Language(id="en"),
        TranslatedText(key="TEST_TEXT_KEY", language="sv", value="sv-text-val"),
        TranslatedText(key="TEST_TEXT_KEY", language="en", value="en-text-val"),
        TranslatedText(key="OTHER_TEXT_KEY", language="sv", value="sv-other-val"),
        TranslatedText(key="OTHER_TEXT_KEY", language="en", value="en-other-val"),
        TranslatedText(key="ONLY_SV_TEXT_KEY", language="sv", value="sv-only-val"),
    ]

    with TestEnvironment(items) as client:
        sv_headers_ok = headers(language="sv")
        res_sv = client.get("/v1/texts/key/TEST_TEXT_KEY", headers=sv_headers_ok)
        assert res_sv.status_code == status.HTTP_200_OK
        assert res_sv.get_json()["TEST_TEXT_KEY"] == "sv-text-val"

        en_headers_ok = headers(language="en")
        res_en = client.get("/v1/texts/key/TEST_TEXT_KEY", headers=en_headers_ok)
        assert res_en.status_code == status.HTTP_200_OK
        assert res_en.get_json()["TEST_TEXT_KEY"] == "en-text-val"

        no_lang_headers = headers(language=None)
        res_no_lang = client.get("/v1/texts/key/TEST_TEXT_KEY", headers=no_lang_headers)
        assert res_no_lang.status_code == status.HTTP_400_BAD_REQUEST

        en_headers = headers(language="en")
        res_missing = client.get("/v1/texts/key/ONLY_SV_TEXT_KEY", headers=en_headers)
        assert res_missing.status_code == status.HTTP_404_NOT_FOUND

        wrong_lang_headers = headers(language="xy")
        res_wrong_lang = client.get(
            "/v1/texts/key/TEST_TEXT_KEY", headers=wrong_lang_headers
        )
        assert res_wrong_lang.status_code == status.HTTP_400_BAD_REQUEST


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
        TextGroup(id="MOBILE_APP"),
        TextGroupMembership(text_key="TEST_TEXT_KEY", group_id="MOBILE_APP"),
        TextGroupMembership(text_key="OTHER_TEXT_KEY", group_id="MOBILE_APP"),
        TextGroupMembership(text_key="ONLY_SV_TEXT_KEY", group_id="MOBILE_APP"),
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
        assert len(body) == 2
        assert body["TEST_TEXT_KEY"] == "en-text-val"
        assert body["OTHER_TEXT_KEY"] == "en-other-val"
        assert "ONLY_SV_TEXT_KEY" not in body
