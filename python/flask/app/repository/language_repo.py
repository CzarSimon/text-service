# Standard libraries
from typing import Optional

# 3rd party libraries
from app.httputil.instrumentation import trace

# Internal modules
from app.models import Language


@trace("language_repo")
def find_language(id: str) -> Optional[Language]:
    return Language.query.filter(Language.id == id).first()

