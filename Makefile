.PHONY: dev
dev:
	docker-compose up --build

.PHONY: clean
clean:
	docker-compose rm
