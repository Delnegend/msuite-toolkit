# mobility-suite-toolkit

## Overview
- **Purpose:** Toolkit for mass import/export and action execution on M-Suite systems (bulk user/device operations, reporting, and automated workflows).

## Prerequisites
- **Windows:** Install [just](https://github.com/casey/just) and [Go](https://go.dev/).
- **Linux / macOS:** Recommended Visual Studio Code with a DevContainer configured for Go, or install `just` and Go locally.

## Setup

### Justfile Configuration
Before using the tool, copy the appropriate file for your operating system to `.justfile` so you don't have to specify the file when running commands:

- **Windows (PowerShell):**
  ```powershell
  Copy-Item .windows.justfile .justfile
  ```
- **Linux / macOS (bash):**
  ```bash
  cp .linux.justfile .justfile
  ```

### Config file (config.toml)
Copy and edit `config.toml` with your environment values:
- `bearer_token`: Open the Admin Portal in your browser, open Developer Tools -> Application -> Local Storage -> select the admin portal origin -> find the key `admin_portal_access_token` and copy its value.
- `admin_user_id`: In the Admin Portal navigate to **Identity > Users, Groups & Unit > Users**, find the currently-logged-in admin user, click the result, then copy the **User ID** from Basic info.
- `admin_portal_address`: Set to the Admin Portal host:port you are using (for example `10.0.0.1:9443`).
- `worker_count`: The number of concurrent workers to use for API requests (default: `100`).

## Usage
After copying the correct `.justfile` for your system, you can build and run using `just`:

- **Build a specific tool:**
  ```bash
  just build NAME
  ```
  (Where `NAME` is one of the folders in `cmd/`, like `users-devices`)

- **Run build-and-test:**
  ```bash
  just run NAME config.toml
  ```

- **Run tests for a specific package:**
  ```bash
  just test PACKAGE
  ```
  (Where `PACKAGE` is a path within `pkg/`, like `endpoints/get-users`)

### Tool Flags
Most tools support the following flags:
- `-config`: Path to config file (default: `./config.toml`)
- `-output`: Path to output CSV file
- `-h` or `-help`: Show help

## Build & Distribute
Use the `build` recipe to produce a user-facing directory in `dist/{{tool-name}}/` which includes the binary, the tool-specific README, and a template `config.toml`. Deliver the entire folder to end users.

## See also
- Repository config template: [config.toml](config.toml)
