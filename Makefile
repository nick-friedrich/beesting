
# Create a new application
new:
	@go run ./cmd/beesting new $(filter-out $@,$(MAKECMDGOALS))

# Run an application in development mode
dev:
	@APP_NAME=$(filter-out $@,$(MAKECMDGOALS)); \
	if [ -f "app/$$APP_NAME/package.json" ]; then \
		cd app/$$APP_NAME && npm run dev; \
	else \
		go run ./cmd/beesting dev $$APP_NAME; \
	fi

# Catch-all target to prevent make errors with app names
%:
	@:

# Show help
help:
	@echo "Beesting Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make new <app-name>  - Create a new application"
	@echo "  make dev <app-name>  - Run an application in dev mode"

