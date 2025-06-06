DOCKER_COMPOSE=docker-compose

GREEN=\033[0;32m
NC=\033[0m

.PHONY: build
build:
	@echo "${GREEN}Building image...${NC}"
	@${DOCKER_COMPOSE} build

.PHONY: rebuild
rebuild:
	@echo "${GREEN}Rebuilding image...${NC}"
	@${DOCKER_COMPOSE} up -d --build

.PHONY: up
up:
	@echo "${GREEN}Starting containers...${NC}"
	@${DOCKER_COMPOSE} --env-file .env up -d

.PHONY: down
down:
	@echo "${GREEN}Stopping containers...${NC}"
	@${DOCKER_COMPOSE} down

.PHONY: restart
restart:
	@echo "${GREEN}Restarting containers...${NC}"
	@${DOCKER_COMPOSE} --env-file .env restart

.PHONY: postgres
postgres:
	@echo "${GREEN}Connecting to database...${NC}"
	@${DOCKER_COMPOSE} exec postgres psql -U postgres -d weather_api


.PHONY: help
help:
	@echo "Available commands:"
	@echo "  ${GREEN}build${NC}       - Build the Docker image"
	@echo "  ${GREEN}rebuild${NC}     - Rebuild the Docker image"
	@echo "  ${GREEN}up${NC}          - Start the containers"
	@echo "  ${GREEN}down${NC}        - Stop the containers"
	@echo "  ${GREEN}restart${NC}     - Restart the containers"
	@echo "  ${GREEN}postgres${NC}    - Connect to the PostgreSQL database"
	@echo "  ${GREEN}help${NC}        - Show this help message"

.DEFAULT_GOAL := help