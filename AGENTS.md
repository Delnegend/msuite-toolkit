# AGENTS.md

> Project knowledge for AI coding agents. Read this before making changes.
> **After making structural changes, update this file** to reflect the current state.
> Add new endpoints to the directory tree, new tools to the tool list, and update patterns that change.

## Tool categorization

| Tool | Type | Config source | Notes |
|------|------|--------------|-------|
| `delete-enrollment-requests` | 🔴 Destructive | Own `config.toml` | OU required, supports `exclude_emails`, `dry_run` |
| `delete-pending-enrollment-requests` | 🔴 Destructive | Root `config.toml` | Only deletes `Pending` status, `dry_run` |
| `set-inactive-users` | 🔴 Destructive | Own `config.toml` | Interactive confirmations, `dry_run` |
| `get-active-users` | 🟢 Read-only | Own `config.toml` | CSV output only |
| `get-provision-policies` | 🟢 Read-only | Root `config.toml` | CSV output only |
| `get-user-devices` | 🟢 Read-only | Root `config.toml` | CSV output only |
| `get-users-history` | 🟢 Read-only | Root `config.toml` | CSV output only |
| `get-users-logins` | 🟢 Read-only | Root `config.toml` | CSV output only |
| `map-apps-to-users` | 🟢 Read-only | Own `config.toml` | CSV output only, supports `filter_by` |
| `playground` | N/A | N/A | Experimental, no README |

🔴 Destructive tools **must** implement `dry_run` with:
- `dry_run = true` → write `to-be-*.csv`, skip mutation, log "Would delete"
- `dry_run = false` → execute mutation, write `deleted-*.csv`, log "Deleted"

## Architecture

```
msuite-toolkit/
├── cmd/                          # Independent CLI tools (one binary per subdirectory)
│   ├── delete-enrollment-requests/
│   ├── delete-pending-enrollment-requests/
│   ├── get-active-users/
│   ├── get-provision-policies/
│   ├── get-user-devices/
│   ├── get-users-history/
│   ├── get-users-logins/
│   ├── map-apps-to-users/
│   ├── playground/
│   └── set-inactive-users/
├── pkg/
│   ├── app/                      # Shared config loading (init.go)
│   ├── endpoints/                # API client functions, one subpackage per endpoint
│   │   ├── delete-enrollment-request/
│   │   ├── find-user-by-email/
│   │   ├── get-access-rules/
│   │   ├── get-devices/
│   │   ├── get-enrollment-requests/
│   │   ├── get-enrollment-requests-by-ouid/   # Composite: gets users by OU, then per-user enrollment requests
│   │   ├── get-provision-policies/
│   │   ├── get-user-apps/
│   │   ├── get-user-device-last-ip/
│   │   ├── get-user-devices/
│   │   ├── get-user-failed-logins/
│   │   ├── get-user-mfa/
│   │   ├── get-users/
│   │   └── inactive-user/
│   ├── httpclient/               # Singleton HTTP client
│   ├── types/                    # Shared types: AppState, UserInfo, DeviceInfo, QueryRequestPayload, etc.
│   └── utils/                    # Utilities (PrintProgressBar)
├── config.toml                   # Default root config (safe template with CHANGE_ME placeholders)
├── config.test.toml              # Config for tests (safe template)
├── config.mytel.toml             # Per-tenant config overrides (real credentials, gitignored)
├── config.viettel-vtn.toml
├── config.vtn-vtn.toml
├── go.mod
└── AGENTS.md
```

### Key architectural rules

1. **Every `cmd/*` is an independent binary.** No code is shared between `cmd/*` packages. Shared helpers go in `pkg/`.
2. **Config is loaded once** via `app.Init(defaultOutput)`, which parses `-config`, `-output`, `-h` flags and decodes TOML into `types.AppState`. The global `app.AppState` is then used via `as := &app.AppState`.
3. **API calls are paginated** via `pond` worker pool. The pool size is `as.WorkerCount` (default 100, set in `pkg/app/init.go`).
4. **Every API endpoint package** exports three layers:
   - `GetX(as, payload)` — fetches one page
   - `GetAllX(as, basePayload, progressChan)` — paginates via pool
   - `GetXWithProgress(as, basePayload)` — prints progress bar wrapper

### `main()` function style — STRICT

**`main()` must be a short pipeline of named function calls.** The following are FORBIDDEN inside `main()`:
- ❌ `for` loops (including `for _, ... := range`)
- ❌ `csv.NewWriter(...)` or any `csvFile.Close()` calls
- ❌ `json.Marshal(...)` calls
- ❌ Anonymous functions (`func() { ... }()`)
- ❌ Raw `os.Create(...)` calls

Helper functions go below `main()` in the same file. Naming conventions:
| Prefix | Purpose | Example |
|--------|---------|---------|
| `fetch*` | Retrieves data from API with progress | `fetchAllEnrollmentRequests` |
| `write*` | Writes a CSV report | `writeOneAppToManyUsersCSV` |
| `compute*` | Filters/transforms data (pure function) | `computeIncludedRequests` |
| `report*` | Writes a report CSV with logging | `reportIncludedRequests` |
| `build*` | Constructs a map/set/index | `buildExcludeUserIDs` |
| `marshal*` | Serializes to JSON/string | `marshalAppsJSON` |

Example — every `main()` should look like this:

```go
func main() {
    outputPath := app.Init("apps_to_users.csv")
    as := &app.AppState

    users := fetchUsers(as)
    appsMap := fetchUserApps(as, users)
    filteredAppsMap := filterAppsByDestination(appsMap, as.FilterBy.DestinationHost, as.FilterBy.DestinationPort)

    writeOneAppToManyUsersCSV(outputPath, filteredAppsMap)
    writeOneUserToOneAppCSV(outputPath, filteredAppsMap)
    writeOneUserToManyAppsCSV(outputPath, filteredAppsMap)
}
```

### Data flow

```
Config (TOML) → types.AppState → endpoint functions → []types.* / map[ID]Data → CSV files
```

All API calls use `Authorization: Bearer {token}` header. Request URLs follow this pattern:

```
https://{admin_portal_address}/enrollment-api/v1/domains/default/enrollment_requests?ctx.user_id={admin_user_id}&request_payload={json}
```

Each API endpoint constructs the URL, marshals the `QueryRequestPayload` as a JSON string into the `request_payload` query parameter.

## Standards

### Go
- Module: `msuite-toolkit`, Go 1.25.6
- Only two external dependencies: `toml` (config parsing), `pond/v2` (worker pool). Everything else is stdlib.
- Package names: snake_case matching the directory name.
  - ⚠️ Directory `find-user-by-email` → package `finduserbyemail` (Go disallows hyphens in package names)
  - Directory `get-enrollment-requests` → package `get_enrollment_requests`
- **Always alias `pkg/endpoints` imports** with snake_case aliases:
  ```go
  get_enrollment_requests "msuite-toolkit/pkg/endpoints/get-enrollment-requests"
  ```
  NEVER use bare imports of endpoint packages.

### When to use `pond` worker pool
- **Always** for paginated API calls (`GetAllX` functions)
- **Yes** for parallelizing independent lookups (e.g. resolving multiple emails to user IDs)
- **Not needed** for single API calls with no parallelism

### CSV writing pattern
Every CSV writer must follow this exact pattern:
```go
csvFile, err := os.Create(fileName)
if err != nil { ... os.Exit(1) }
defer func() { csvFile.Close() }()

w := csv.NewWriter(csvFile)
w.Comma = '|'            // pipe-delimited
defer w.Flush()

// header → rows
if err := w.Write(header); err != nil { ... }
for _, item := range items {
    if err := w.Write(row); err != nil { ... }
}

// mandatory final check
if err := w.Error(); err != nil { ... os.Exit(1) }
```

### Filter format
Server-side filters in `QueryRequestPayload` use this structure:
```go
[]any{map[string]any{
    "key":      "field_name",      // e.g. "status", "user_id", "device_id"
    "operator": "equal_to",        // always "equal_to"
    "value":    "field_value",
}}
```

### Error handling
- `slog.Error(msg, "err", err)` then `os.Exit(1)` — no error returns from `main()`
- Worker pool tasks: log errors internally, collect them, return aggregated error
- Non-2xx HTTP: log status + body, return as error

## Testing

- `*_test.go` files live alongside source (both `cmd/*` and `pkg/*`)
- Unit tests for pure functions (`computeIncludedRequests`, `filterAppsByDestination`)
- Integration tests load `../../../config.test.toml` and call real API endpoints
- Integration tests **skip** when required config fields are empty:
  ```go
  if appState.OrganizationalUnitID == "" {
      t.Skip("OrganizationalUnitID is empty in config.test.toml, skipping test")
  }
  ```

## Security

- **Real bearer tokens are sensitive.** Never commit them. All config templates use `"CHANGE_ME"`.
- Per-tenant files (`config.mytel.toml`, etc.) contain real tokens — they must be `.gitignore`d.
- `config.toml` and `config.test.toml` are safe templates with placeholder tokens only.
- All HTTP requests use `Authorization: Bearer {token}` — no other auth mechanism.
- HTTP client from `pkg/httpclient.GetHTTPClient()` — singleton, no custom TLS.

## Deployment

### Build
```bash
go build ./cmd/...                              # All tools
go build ./cmd/delete-enrollment-requests/      # Single tool
```

### Run
```bash
./delete-enrollment-requests.exe -config ./config.toml
```

### Flags (standard across all tools via `app.Init`)
| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `./config.toml` | Path to TOML config |
| `-output` | tool-specific | Output CSV file path |
| `-h`, `-help` | — | Show help text |

### Documentation conventions
- Every `cmd/*` has `README.md` (English) and `README.vi.md` (Vietnamese)
- README sections: Description, Quick steps, Run notes, Flags, config.toml how-to, Reports (CSV), Example
- Per-tool `config.toml` files are templates with `CHANGE_ME` placeholders
- Binary names use `.exe` in README examples (Windows toolchain)
