# Use sh on Linux/macOS
set shell := ["sh", "-cu"]

# Build a CLI tool from ./cmd/<name> into ./dist/<name>
build NAME:
	rm -rf "dist/{{NAME}}"
	mkdir -p "dist/{{NAME}}"
	go build -o "dist/{{NAME}}/{{NAME}}" "./cmd/{{NAME}}"
	[ -f "cmd/{{NAME}}/README.md" ] && cp "cmd/{{NAME}}/README.md" "dist/{{NAME}}/README.md" || true
	[ -f "cmd/{{NAME}}/README.vi.md" ] && cp "cmd/{{NAME}}/README.vi.md" "dist/{{NAME}}/README.vi.md" || true
	# Prefer a command-specific config if present, otherwise fall back to repo root config.toml
	if [ -f "cmd/{{NAME}}/config.toml" ]; then \
		cp "cmd/{{NAME}}/config.toml" "dist/{{NAME}}/config.toml"; \
	elif [ -f "config.toml" ]; then \
		cp "config.toml" "dist/{{NAME}}/config.toml"; \
	fi
	echo "Built dist/{{NAME}}/{{NAME}} with documentation and config"

# Run a built CLI with a config file and optional output file: run <name> <path-to-config> [output-file]
run NAME CONFIG OUTPUT="":
	go run ./cmd/{{NAME}} -config {{CONFIG}} {{if OUTPUT == "" { "" } else { "-output " + OUTPUT }}}

# Run tests for a specific package: test <package-path>
test PACKAGE:
	go test -v "./pkg/{{PACKAGE}}"

# Check: compile all CLI packages under cmd/ to ensure they build
check:
	@echo "Checking CLI compilation..."
	go build ./cmd/...
	@echo "OK: CLI packages compile"

format:
	go fmt ./...