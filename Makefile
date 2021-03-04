run:
	@revel run

build:
	@docker-compose build

up:
	@docker-compose up -d

logs:
	docker-compose logs -f

ps:
	@docker-compose ps

stop:
	@docker-compose stop

rm: stop
	@docker-compose rm
