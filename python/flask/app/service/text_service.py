# Standard library
from typing import Dict, Optional

# 3rd party modules
from app.httputil.error import NotFoundError, BadRequestError
from app.httputil.instrumentation import trace

# Internal modules
from app.models import TranslatedText, Language, TextGroup
from app.repository import text_repo, language_repo


@trace("text_service")
def get_text_by_key(key: str, language: str) -> Dict[str, str]:
    _assert_language_support(language)
    text = text_repo.find_by_key(key, language)
    if not text:
        raise NotFoundError()
    return {text.key: text.value}


@trace("text_service")
def get_text_by_group(group_id: str, language: str) -> Dict[str, Optional[str]]:
    _assert_language_support(language)
    group = _assert_group(group_id)
    keys = [text_group.text_key for text_group in group.texts]
    texts = text_repo.find_by_keys(keys, language)
    return {text.key: text.value for text in texts}


@trace("text_service")
def _assert_group(group_id: str) -> TextGroup:
    group = text_repo.find_group(group_id)
    if not group:
        raise NotFoundError(f"No such group: {group_id}")
    return group


@trace("text_service")
def _assert_language_support(language_id: str) -> Language:
    lang = language_repo.find_language(language_id)
    if not lang:
        raise BadRequestError(f"Unsupported language: {language_id}")
    return lang
