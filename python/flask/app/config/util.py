# Standard library
import os


def get_dsn(storage_type: str) -> str:
    if storage_type == "memory":
        return "sqlite:///:memory:"
    if storage_type == "postgres":
        return _get_postgres_dsn()
    if storage_type == "sqlite":
        return _get_sqlite_dsn()
    return _get_postgres_dsn()


def _get_postgres_dsn() -> str:
    DB_USERNAME = os.environ["DB_USERNAME"]
    DB_PASSWORD = os.environ["DB_PASSWORD"]
    DB_HOST = os.environ["DB_HOST"]
    DB_PORT = os.environ["DB_PORT"]
    DB_NAME = os.environ["DB_NAME"]
    return f"postgresql://{DB_USERNAME}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}"


def _get_sqlite_dsn() -> str:
    return os.environ["DB_NAME"]
