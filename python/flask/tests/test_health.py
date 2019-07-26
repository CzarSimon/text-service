# Standard library
from unittest.mock import patch

# 3rd party modules
from app.httputil import status

# Intenal modules
from tests import TestEnvironment, JSON


def test_get_health_success():
    with TestEnvironment() as client:
        res = client.get("/health", content_type=JSON)
        assert res.status_code == status.HTTP_200_OK
        assert res.get_json()["status"] == "UP"
        assert res.get_json()["db"] == "UP"


def test_get_health_fail():
    with patch("app.db.engine.execute", side_effect=Exception("mock fail")):
        with TestEnvironment() as client:
            res = client.get("/health", content_type=JSON)
            assert res.status_code == status.HTTP_503_SERVICE_UNAVAILIBLE
