# Use sh on Linux/macOS
set shell := ["sh", "-cu"]

# Build a CLI tool from ./cmd/<name> into ./dist/<name>
build NAME:
	mkdir -p "dist/{{NAME}}"
	go build -o "dist/{{NAME}}/{{NAME}}" "./cmd/{{NAME}}"
	[ -f "cmd/{{NAME}}/README.md" ] && cp "cmd/{{NAME}}/README.md" "dist/{{NAME}}/README.md" || true
	[ -f "config.toml" ] && cp "config.toml" "dist/{{NAME}}/config.toml" || true
	echo "Built dist/{{NAME}}/{{NAME}} with documentation and config"

# Run a built CLI with a config file: run <name> <path-to-config>
run NAME CONFIG:
	if [ ! -f "dist/{{NAME}}/{{NAME}}" ]; then echo "Executable dist/{{NAME}}/{{NAME}} not found, run 'just build {{NAME}}' first"; exit 1; fi
	./dist/{{NAME}}/{{NAME}} -config {{CONFIG}}
