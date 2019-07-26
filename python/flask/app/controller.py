# Standard library
import logging

# 3rd party modules
import flask
from flask import jsonify, make_response, request
from app import httputil
from app.httputil import status
from app.httputil.error import BadRequestError
from app.httputil.instrumentation import trace

# Internal modules
from app.config import LANGUAGE_HEADER
from app.service import health


_log = logging.getLogger(__name__)


@trace("controller")
def get_text_by_key(key: str) -> flask.Response:
    return httputil.create_ok_response()


@trace("controller")
def get_text_group(group_id: str) -> flask.Response:
    return httputil.create_ok_response()


@trace("controller")
def check_health() -> flask.Response:
    health_status = health.check()
    return httputil.create_response(health_status)


def _get_language() -> str:
    lang = request.headers.get(LANGUAGE_HEADER)
    if not lang:
        raise BadRequestError("No language specified")
    return lang
