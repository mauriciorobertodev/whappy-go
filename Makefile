dev:
	clear && go run ./cmd/api/main.go

build:
	clear && go build -o ./bin/api ./cmd/api/main.go

test:
	clear && ginkgo ./...

test-bus:
	clear && ginkgo -v ./internal/infra/eventbus/...

test-cache:
	clear && ginkgo -v ./internal/infra/cache/...

test-storage:
	clear && ginkgo -v ./internal/infra/storage/...

test-repository:
	clear && ginkgo -v ./internal/infra/repository/...

test-webhook:
	clear && ginkgo -v ./internal/domain/webhook/...

docker-test-up:
	clear && docker-compose -f ./docker/docker-compose.test.yml --env-file ./.env.test up -d --remove-orphans

docker-test-down:
	clear && docker-compose -f ./docker/docker-compose.test.yml down -v --remove-orphans