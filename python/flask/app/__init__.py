# Standard library
import logging
from datetime import datetime
from functools import wraps
from typing import Callable, Optional
from uuid import uuid4

# 3rd party modules.
import flask
import werkzeug
from flask import Flask, jsonify, make_response, request
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
from app.httputil import status, instrumentation
from app.httputil.instrumentation import setup_instrumentation
from app.httputil.error import (
    RequestError,
    NotFoundError,
    MethodNotAllowedError,
    InternalServerError,
)

# Internal modules
from app.config import AppConfig
from app.config import REQUEST_ID_HEADER, SERVICE_NAME, SERVICE_VERSION


app = Flask(SERVICE_NAME)
app.config.from_object(AppConfig)
db = SQLAlchemy(app)
migrate = Migrate(app, db)
setup_instrumentation(app, name=SERVICE_NAME, version=SERVICE_VERSION)


from app import routes

# from app import models


error_log = logging.getLogger("ErrorHandler")


@app.errorhandler(RequestError)
def handle_request_error(error: RequestError) -> flask.Response:
    """Handles errors encountered when handling requests.

    :param error: Encountered RequestError.
    :return: flask.Response indicating the encountered error.
    """
    if error.status() >= status.HTTP_500_INTERNAL_SERVER_ERROR:
        error_log.error(str(error))
    else:
        error_log.info(str(error))
    json_error = jsonify(error.asdict())
    return make_response(json_error, error.status())


@app.errorhandler(404)
def handle_not_found(err: werkzeug.exceptions.NotFound) -> flask.Response:
    """Handles 404 errors.

    :return: flask.Response indicating the encountered error.
    """
    error = NotFoundError()
    return make_response(jsonify(error.asdict()), error.status())


@app.errorhandler(405)
def handle_method_not_allowed(
    err: werkzeug.exceptions.MethodNotAllowed
) -> flask.Response:
    """Handles 405 errors.

    :return: flask.Response indicating the encountered error.
    """
    error = MethodNotAllowedError()
    return make_response(jsonify(error.asdict()), error.status())


@app.errorhandler(Exception)
def handle_internal_error(err: Exception) -> flask.Response:
    """Handles uncaught internal errors, logs and sends as internal errors.

    :return: flask.Response indicating the encountered error.
    """
    request_id = instrumentation.get_request_id()
    error_log.error(f"Error: [{err}] requestId={request_id}")
    error = InternalServerError()
    return make_response(jsonify(error.asdict()), error.status())
