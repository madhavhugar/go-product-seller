init:
	bash build/install_server.sh

db-up:
	docker-compose -f build/docker-compose.yml start db

up:
	docker-compose -f build/docker-compose.yml up

stop:
	docker-compose -f build/docker-compose.yml stop
