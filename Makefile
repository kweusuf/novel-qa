.PHONY: help build run stop clean dev build-prod run-prod logs shell

# Auto-detect container runtime
DOCKER := $(shell command -v docker 2> /dev/null || command -v podman 2> /dev/null)
COMPOSE := $(shell command -v docker-compose 2> /dev/null || command -v podman-compose 2> /dev/null)

# Fallback detection if the above doesn't work
ifeq ($(DOCKER),)
    DOCKER := podman
endif
ifeq ($(COMPOSE),)
    COMPOSE := podman-compose
endif

# Show which runtime we're using
RUNTIME_INFO := $(shell basename $(DOCKER))
COMPOSE_INFO := $(shell basename $(COMPOSE))

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "🔧 Using runtime: $(RUNTIME_INFO)"
	@echo "🔧 Using compose: $(COMPOSE_INFO)"
	@echo ""
	@echo "💡 To force a specific runtime:"
	@echo "   make DOCKER=docker COMPOSE=docker-compose build"
	@echo "   make DOCKER=podman COMPOSE=podman-compose build"

build: ## Build the production container image
	@echo "🔨 Building with $(RUNTIME_INFO)..."
	$(DOCKER) build -t novel-qa:latest .

build-dev: ## Build the development container image
	@echo "🔨 Building development image with $(RUNTIME_INFO)..."
	$(DOCKER) build -f Dockerfile.dev -t novel-qa:dev .

run: ## Run the production application
	@echo "🚀 Starting production services with $(COMPOSE_INFO)..."
	$(COMPOSE) up -d

run-dev: ## Run the development application with hot reloading
	@echo "🚀 Starting development services with $(COMPOSE_INFO)..."
	$(COMPOSE) -f docker-compose.dev.yml up -d

stop: ## Stop all containers
	@echo "🛑 Stopping all services..."
	$(COMPOSE) down
	$(COMPOSE) -f docker-compose.dev.yml down

stop-prod: ## Stop production containers
	@echo "🛑 Stopping production services..."
	$(COMPOSE) down

stop-dev: ## Stop development containers
	@echo "🛑 Stopping development services..."
	$(COMPOSE) -f docker-compose.dev.yml down

clean: ## Remove all containers, images, and volumes
	@echo "🧹 Cleaning up everything..."
	$(COMPOSE) down -v --rmi all
	$(COMPOSE) -f docker-compose.dev.yml down -v --rmi all
	$(DOCKER) system prune -f

logs: ## View logs from all services
	$(COMPOSE) logs -f

logs-dev: ## View logs from development services
	$(COMPOSE) -f docker-compose.dev.yml logs -f

shell: ## Open shell in the running container
	$(COMPOSE) exec novel-qa sh

shell-dev: ## Open shell in the development container
	$(COMPOSE) -f docker-compose.dev.yml exec novel-qa sh

pull-ollama: ## Pull Ollama models (run after starting services)
	@echo "📥 Pulling Ollama models..."
	$(COMPOSE) exec ollama ollama pull phi3
	$(COMPOSE) exec ollama ollama pull llama3
	$(COMPOSE) exec ollama ollama pull mistral
	$(COMPOSE) exec ollama ollama pull gemma

status: ## Show status of all containers
	@echo "📊 Container status:"
	$(COMPOSE) ps
	@echo ""
	@echo "📊 Development container status:"
	$(COMPOSE) -f docker-compose.dev.yml ps

restart: ## Restart all services
	@echo "🔄 Restarting all services..."
	$(COMPOSE) restart

restart-dev: ## Restart development services
	@echo "🔄 Restarting development services..."
	$(COMPOSE) -f docker-compose.dev.yml restart

# Force specific runtime commands
docker: ## Force use of Docker runtime
	$(MAKE) DOCKER=docker COMPOSE=docker-compose $(MAKE_ARGS)

podman: ## Force use of Podman runtime
	$(MAKE) DOCKER=podman COMPOSE=podman-compose $(MAKE_ARGS)

# Show runtime info
runtime-info: ## Show which runtime is being used
	@echo "🔧 Container Runtime: $(RUNTIME_INFO)"
	@echo "🔧 Compose Tool: $(COMPOSE_INFO)"
	@echo "🔧 Full Docker Path: $(DOCKER)"
	@echo "🔧 Full Compose Path: $(COMPOSE)"
	@echo ""
	@echo "💡 Available runtimes:"
	@echo "   Docker: $(shell command -v docker 2> /dev/null || echo "Not found")"
	@echo "   Docker Compose: $(shell command -v docker-compose 2> /dev/null || echo "Not found")"
	@echo "   Podman: $(shell command -v podman 2> /dev/null || echo "Not found")"
	@echo "   Podman Compose: $(shell command -v podman-compose 2> /dev/null || echo "Not found")"

