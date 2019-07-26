# 3rd party modules.
import flask

# Internal modules
from app import app
from app import controller
from app import httputil


@app.route("/v1/texts/key/<key>", methods=["GET"])
def get_text(key: str) -> flask.Response:
    return controller.get_text_by_key(key)


@app.route("/v1/texts/group/<group_id>", methods=["GET"])
def get_text_group(group_id: str) -> flask.Response:
    return controller.get_text_group(group_id)


@app.route("/health", methods=["GET"])
def check_health() -> flask.Response:
    return controller.check_health()

