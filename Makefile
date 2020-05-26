.PHONY: dev
dev: SHELL:=/bin/bash
dev:
	source .environ && docker-compose up --build

.PHONY: clean
clean:
	docker-compose rm
