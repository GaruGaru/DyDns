VERSION=$(shell git rev-parse --short HEAD)
DIST_DIR=./dist

# GO

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	go build -o ${DIST_DIR}/bets bets.go

.PHONY: test
test:
	go test -v ./...

.PHONY: deps
deps:
	dep ensure


# Docker

IMAGE=garugaru/dydns
DEV_COMPOSE=docker/docker-compose-dev.yml
PROD_COMPOSE=docker/docker-compose-prod.yml

.PHONY: docker-push
docker-push: docker-build-prod
	docker push garugaru/dydns

.PHONY: docker-build-dev
docker-build-dev:
	docker-compose -f ${DEV_COMPOSE} build


.PHONY: docker-build-prod
docker-build-prod:
	docker-compose -f ${PROD_COMPOSE} build


.PHONY: docker-up
docker-up:
	docker-compose -f ${DEV_COMPOSE} up


.PHONY: docker-down
docker-down:
	docker-compose -f ${DEV_COMPOSE} down
