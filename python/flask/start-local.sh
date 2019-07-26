export prometheus_multiproc_dir=/tmp/text-service/prom-data
export FLASK_APP=run.py

flask db upgrade
flask run --port 8080
