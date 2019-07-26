export DB_USERNAME=textservice
export DB_PASSWORD=password
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_NAME=texts


test:
	sh run-tests.sh

run-local:
	sh start-local.sh

image:
	sh build-image.sh

db-init:
	rm -rf migrations
	flask db init

migration:
	flask db migrate

db-upgrade:
	flask db upgrade

db-downgrade:
	flask db downgrade