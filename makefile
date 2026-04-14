MAIN_FILE	=main.go
BINARY_NAME =go_forum_backend
BUILD_DIR	=build

serve:
	@go run $(MAIN_FILE) serve

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Compilation de $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Binaire généré : $(BUILD_DIR)/$(BINARY_NAME)"

migrate:
	@go run $(MAIN_FILE) migrate

generate-keys:
	@echo "Generate Private Key :"
	@openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
	@echo "Generate Public Key :"
	@openssl pkey -in private_key.pem -pubout -out public_key.pem
	@echo "Keys Generate Successfully !!!"

clean:
	@echo "Nettoyage des fichiers générés..."
	@rm -rf $(BUILD_DIR)
