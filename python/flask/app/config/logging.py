# Standard library
import logging


LOGGING_CONIFG = {
    "version": 1,
    "formatters": {
        "default": {
            "format": "%(asctime)s %(name)s %(levelname)s %(message)s %(lineno)s",
            "class": "pythonjsonlogger.jsonlogger.JsonFormatter",
        }
    },
    "handlers": {
        "console": {
            "class": "logging.StreamHandler",
            "formatter": "default",
            "level": logging.DEBUG,
        }
    },
    "root": {"handlers": ["console"], "level": logging.DEBUG},
}

logging.getLogger("werkzeug").setLevel(logging.ERROR)
