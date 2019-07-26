# Standard libraries
from typing import List, Optional

# 3rd party libraries
from app.httputil.instrumentation import trace

# Internal modules
from app.models import TranslatedText, TextGroup, Group


@trace("text_repo")
def find_by_key(key: str, language: str) -> Optional[TranslatedText]:
    return TranslatedText.query.filter(
        TranslatedText.key == key, TranslatedText.language == language
    ).first()


@trace("text_repo")
def find_by_keys(keys: List[str], language: str) -> List[TranslatedText]:
    return TranslatedText.query.filter(
        TranslatedText.key.in_(keys),  # type: ignore
        TranslatedText.language == language,
    ).all()


@trace("text_repo")
def find_group(group_id) -> Optional[Group]:
    return Group.query.filter(Group.id == group_id).first()
