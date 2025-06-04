PHONY: run restart

run: restart
	docker-compose exec app /bin/sh

restart:
	docker-compose down
	docker-compose up -d
	