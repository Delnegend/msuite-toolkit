# Use PowerShell on Windows
set shell := ["powershell.exe", "-NoLogo", "-Command"]

# Build a CLI tool from ./cmd/<name> into ./dist/<name>
build NAME:
	if (!(Test-Path -Path "dist/{{NAME}}")) { New-Item -ItemType Directory -Path "dist/{{NAME}}" -Force | Out-Null }
	go build -o "dist/{{NAME}}/{{NAME}}.exe" "./cmd/{{NAME}}"
	if (Test-Path -Path "cmd/{{NAME}}/README.md") { Copy-Item "cmd/{{NAME}}/README.md" "dist/{{NAME}}/README.md" -Force }
	# Prefer a command-specific config if present, otherwise fall back to repo root config.toml
	if (Test-Path -Path "cmd/{{NAME}}/config.toml") {
		Copy-Item "cmd/{{NAME}}/config.toml" "dist/{{NAME}}/config.toml" -Force
	} elseif (Test-Path -Path "config.toml") {
		Copy-Item "config.toml" "dist/{{NAME}}/config.toml" -Force
	}
	Write-Host "Built dist/{{NAME}}/{{NAME}}.exe with documentation and config"

# Run a built CLI with a config file and optional output file: run <name> <path-to-config> [output-file]
run NAME CONFIG OUTPUT="":
	go run ./cmd/{{NAME}} -config {{CONFIG}} {{if OUTPUT == "" { "" } else { "-output " + OUTPUT }}}

# Run tests for a specific package: test <package-path>
test PACKAGE:
	go test -v "./pkg/{{PACKAGE}}"

# Check: compile all CLI packages under cmd/ to ensure they build
check:
	go build ./cmd/...
	Write-Host "OK: CLI packages compile"

format:
	go fmt ./...