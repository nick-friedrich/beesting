# Example API Makefile

# Install dependencies
install:
	npm install

# Build CSS once
css:
	npm run build-css

# Watch CSS changes
css-watch:
	npm run watch-css

# Development mode (CSS watch + Go dev server)
dev:
	npm run dev

# Help
help:
	@echo "Available commands:"
	@echo "  make install    - Install npm dependencies"
	@echo "  make css       - Build CSS once"
	@echo "  make css-watch - Watch CSS changes"
	@echo "  make dev       - Start development mode (CSS watch + Go dev server)"
	@echo "  make help      - Show this help"

.PHONY: install css css-watch dev help
