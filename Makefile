export USERID=$(shell id -u):$(shell id -g)

.PHONY:
build:
	docker compose build compiler
	docker compose up -d compiler
	docker compose exec compiler go build ./cmd/qre
	mv src/qre .

.PHONY: install
install:
	mkdir -pv $$HOME/.qre
	mkdir -pv $$HOME/.local/bin
	mv qre $$HOME/.local/bin
	$(shell cp -r qres/* $$HOME/.qre)
