# Standard library
import os
from logging.config import dictConfig

# Internal modules
from .logging import LOGGING_CONIFG
from .util import get_dsn

dictConfig(LOGGING_CONIFG)


SERVICE_NAME: str = "text-service"
SERVICE_VERSION: str = "0.2.0"
SERVER_NAME: str = f"{SERVICE_NAME}/{SERVICE_VERSION}"
REQUEST_ID_HEADER: str = "X-Request-ID"
LANGUAGE_HEADER = "Accept-Language"
STORAGE_TYPE: str = os.getenv("STORAGE", "postgres")


class AppConfig:
    SQLALCHEMY_DATABASE_URI: str = get_dsn(STORAGE_TYPE)
    SQLALCHEMY_TRACK_MODIFICATIONS: bool = False

