MAIN_FILE				=main.go
BINARY_NAME 			=go_forum_backend
BUILD_DIR				=build
DOCKER_COMPOSE_DEV_FILE	=docker-compose.dev.yml
APP_NAME				=app_forum

help:
	@echo "Available commands:"
	@echo ""
	@echo "  serve                 - Start the Go server locally"
	@echo "  build                 - Compile the project"
	@echo "  clean                 - Remove generated files"
	@echo ""
	@echo "  docker-dev            - Start containers in dev mode"
	@echo "  docker-dev-up         - Start containers in dev mode"
	@echo "  docker-dev-down       - Stop containers"
	@echo "  docker-dev-logs       - Show container logs"
	@echo ""
	@echo "  migrate               - Run migrations (local)"
	@echo "  docker-migrate        - Run migrations (inside container)"
	@echo ""
	@echo "  generate-keys         - Generate RSA keys (private_key.pem, public_key.pem)"

# --- Development ---
serve:
	@go run $(MAIN_FILE) serve

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Compilation de $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Binaire généré : $(BUILD_DIR)/$(BINARY_NAME)"

clean:
	@echo "Nettoyage des fichiers générés..."
	@rm -rf $(BUILD_DIR)

migrate:
	@go run $(MAIN_FILE) migrate

# --- Docker ---
docker-dev:
	@docker compose -f $(DOCKER_COMPOSE_DEV_FILE) up -d --build

docker-dev-up:
	@docker compose -f $(DOCKER_COMPOSE_DEV_FILE) up -d --build

docker-dev-down:
	@docker compose -f $(DOCKER_COMPOSE_DEV_FILE) down

docker-dev-logs:
	@docker compose -f $(DOCKER_COMPOSE_DEV_FILE) logs -f

docker-migrate:
	@docker exec -it $(APP_NAME) go run $(MAIN_FILE) migrate

# --- Utilities ---
generate-keys:
	@echo "Generate Private Key:"
	@openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
	@echo "Generate Public Key:"
	@openssl pkey -in private_key.pem -pubout -out public_key.pem
	@echo "Keys generated successfully!"
