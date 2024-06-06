# Makefile

# Variables
REGISTRY=registry.gitlab.com/docshade
IMAGES=frontend py-anonymizer notification-service queue-service document-upload-service

# Default rule
all: push_images_into_gitlab

# Rule for tagging and pushing images into GitLab
push_images_into_gitlab: tag_images push_images

tag_images:
	@echo "Tagging images..."
	$(foreach IMAGE, $(IMAGES), docker tag dev-$(IMAGE):latest $(REGISTRY)/$(IMAGE):latest;)

push_images:
	@echo "Pushing images..."
	$(foreach IMAGE, $(IMAGES), docker push $(REGISTRY)/$(IMAGE):latest;)


.PHONY: start_app_dev
start_app_dev:
	@echo "Запускаем приложение..."
	BACKEND_HOST=localhost CONFIG_FILE=config.dev.yaml PLATFORM=linux/arm64 docker-compose -f build/dev/docker-compose.app.yaml up -d

.PHONY: build_images_prod
build_images_prod:
	@echo "Запускаем приложение..."
	BACKEND_HOST=_ CONFIG_FILE=config.prod.yaml PLATFORM=linux/amd64 docker-compose -f build/dev/docker-compose.app.yaml build


# Остановить приложения и базы данных
.PHONY: stop
stop:
	@echo "Останавливаем приложение и базу данных..."
	docker-compose -f build/dev/db/docker-compose.db.yaml down
	docker-compose -f build/dev/docker-compose.app.yaml down

# Запустить тесты
.PHONY: run_tests
run_tests:
	go test rest-executor/usecases/rest_service
	go test rest-executor/entrypoints/http/v1/health

