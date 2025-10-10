
# Create a new application
new:
	@go run ./cmd/beesting new $(filter-out $@,$(MAKECMDGOALS))

# Run an application in development mode
dev:
	@go run ./cmd/beesting dev $(filter-out $@,$(MAKECMDGOALS))

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

