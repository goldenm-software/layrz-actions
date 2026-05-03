.PHONY: help push push-minor push-major current-version build-tools clean-tools

TOOLS := changelog coverage-comment
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

# Default target - show help
help:
	@echo "GitHub Actions CI Workflows - Version Management"
	@echo ""
	@echo "Available targets:"
	@echo "  make push         - Create a new patch version (e.g., v1.0.1 → v1.0.2)"
	@echo "  make push-minor   - Create a new minor version (e.g., v1.0.2 → v1.1.0)"
	@echo "  make push-major   - Create a new major version (e.g., v1.2.3 → v2.0.0)"
	@echo "  make current-version - Show current version"
	@echo ""
	@echo "Examples:"
	@echo "  make push               # v1.0.1 → v1.0.2 (patch - bug fixes)"
	@echo "  make push-minor         # v1.0.2 → v1.1.0 (minor - new features)"
	@echo "  make push-major         # v1.2.3 → v2.0.0 (major - breaking changes)"

# Show current version
current-version:
	@./scripts/current-version.sh

# Create a new patch version and push
push:
	@./scripts/version-push.sh

# Create a new minor version and push
push-minor:
	@./scripts/version-push-minor.sh

# Create a new major version and push
push-major:
	@./scripts/version-push-major.sh

# Build Go tool binaries for all platforms
build-tools:
	@mkdir -p tools/bin
	@$(foreach platform,$(PLATFORMS), \
		$(eval OS=$(word 1,$(subst /, ,$(platform)))) \
		$(eval ARCH=$(word 2,$(subst /, ,$(platform)))) \
		$(foreach tool,$(TOOLS), \
			echo "Building $(tool)-$(OS)-$(ARCH)..." && \
			GOOS=$(OS) GOARCH=$(ARCH) go build -C tools -o bin/$(tool)-$(OS)-$(ARCH) ./$(tool) && \
		) \
	) true

# Remove built binaries
clean-tools:
	@rm -rf tools/bin
