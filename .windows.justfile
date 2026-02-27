# Use PowerShell on Windows
set shell := ["powershell.exe", "-NoLogo", "-Command"]

# Build a CLI tool from ./cmd/<name> into ./dist/<name>
build NAME:
	if (!(Test-Path -Path "dist/{{NAME}}")) { New-Item -ItemType Directory -Path "dist/{{NAME}}" -Force | Out-Null }
	go build -o "dist/{{NAME}}/{{NAME}}.exe" "./cmd/{{NAME}}"
	if (Test-Path -Path "cmd/{{NAME}}/README.md") { Copy-Item "cmd/{{NAME}}/README.md" "dist/{{NAME}}/README.md" -Force }
	if (Test-Path -Path "config.toml") { Copy-Item "config.toml" "dist/{{NAME}}/config.toml" -Force }
	Write-Host "Built dist/{{NAME}}/{{NAME}}.exe with documentation and config"

# Run a built CLI with a config file: run <name> <path-to-config>
run NAME CONFIG:
	if (!(Test-Path -Path "dist/{{NAME}}/{{NAME}}.exe")) { Write-Host "Executable dist/{{NAME}}/{{NAME}}.exe not found, run `just build {{NAME}}` first"; exit 1 }
	./dist/{{NAME}}/{{NAME}}.exe -config {{CONFIG}}
