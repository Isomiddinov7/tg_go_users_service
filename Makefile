CURRENT_DIR=$(shell pwd)

APP=$(shell basename "${CURRENT_DIR}")
APP_CMD_DIR=${CURRENT_DIR}/cmd

TAG=latest
ENV_TAG=latest

pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule update --remote --merge

copy-proto-module:
	rm -rf "${CURRENT_DIR}/protos"
	rsync -rv --exclude={'/.git','LICENSE','README.md'} "${CURRENT_DIR}/tg_protos"/* "${CURRENT_DIR}/protos"

gen-proto-module:
	./scripts/gen_proto.sh "${CURRENT_DIR}"

migration-up:
	migrate -path ./migrations/postgres -database 'postgres://bahodir:1100@0.0.0.0:5432/telegram_bot?sslmode=disable' up

migration-down:
	migrate -path ./migrations/postgres -database 'postgres://bahodir:1100@0.0.0.0:5432/telegram_bot?sslmode=disable' down

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o "${CURRENT_DIR}/bin/${APP}" "${APP_CMD_DIR}/main.go"

run:
	go run cmd/main.go

lint: ## Run golangci-lint with printing to stdout
	golangci-lint -c .golangci.yaml run --build-tags "musl" ./...


git_clear:
	git gc --aggressive && git fetch -p && git rm -r --cached . && git rm --cached crm_protos