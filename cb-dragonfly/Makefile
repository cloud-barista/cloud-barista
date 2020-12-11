
compose-build:
	@echo "Starting Docker Image Build..."
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker build . -t cloudbaristaorg/cb-dragonfly:espresso-v0.1-kafka

compose-up:
	@echo "Starting Docker Compose..."
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker-compose up -d --build

compose-up-dev:
	@echo "Starting Docker Compose..."
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker-compose -f docker-compose-dev-df.yaml up -d --build

compose-rm:
	@echo "Stopping Docker Compose..."
	docker-compose stop && docker-compose rm

clean:
	@echo "Clean Container in cb-dragonfly module..."
	docker ps -aqf name="^cb-dragonfly" | xargs -I {} docker rm -f {}
