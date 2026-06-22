# AGENTS.md

> Read before making changes. Update when adding tools, endpoints, or patterns.

## Tool catalog

| Tool | Type | Config source | Notes |
|------|------|--------------|-------|
| `add-users-to-group` | 🔴 Destructive | Own `config.toml` | Resolves emails → user IDs, adds in batches of 10, `dry_run` |
| `delete-enrollment-requests` | 🔴 Destructive | Own `config.toml` | OU required, deletes only desktop requests, `exclude_emails`, `dry_run` |
| `delete-pending-enrollment-requests` | 🔴 Destructive | Root `config.toml` | Only deletes `Pending` status, `dry_run` |
| `set-inactive-users` | 🔴 Destructive | Own `config.toml` | Locks users via `PUT /users/{id}/lock`. **Typed-text confirmations** (not y/n) if threshold < 3 months, `dry_run` |
| `get-active-users` | 🟢 Read-only | Own `config.toml` | CSV output |
| `get-provision-policies` | 🟢 Read-only | Root `config.toml` | CSV output |
| `get-user-devices` | 🟢 Read-only | Root `config.toml` | CSV output |
| `get-users-history` | 🟢 Read-only | Root `config.toml` | Full history: MFA + devices + last IPs + failed logins |
| `get-users-logins` | 🟢 Read-only | Root `config.toml` | MFA status + failed logins |
| `map-apps-to-users` | 🟢 Read-only | Own `config.toml` | 3 CSV views, supports `filter_by` |
| `playground` | N/A | N/A | Experimental, no README |

🔴 Destructive tools **must** implement `dry_run`:
- `dry_run = true` → write `to-be-*.csv`, skip mutation, log "Would delete"
- `dry_run = false` → execute mutation, write `deleted-*.csv`, log "Deleted"

## Architecture

```
msuite-toolkit/
├── cmd/                      # One independent binary per subdirectory
│   ├── add-users-to-group/
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
│   ├── app/                  # Config loading, flag parsing via Init()
│   ├── endpoints/            # API client packages (one per endpoint)
│   │   ├── add-users-to-group/
│   │   ├── delete-enrollment-request/
│   │   ├── find-user-by-email/
│   │   ├── get-access-rules/
│   │   ├── get-devices/
│   │   ├── get-enrollment-requests/
│   │   ├── get-enrollment-requests-by-ouid/
│   │   ├── get-provision-policies/
│   │   ├── get-user-apps/
│   │   ├── get-user-device-last-ip/
│   │   ├── get-user-devices/
│   │   ├── get-user-failed-logins/
│   │   ├── get-user-mfa/
│   │   ├── get-users/
│   │   └── inactive-user/
│   ├── httpclient/           # Singleton HTTP client (InsecureSkipVerify, 30s timeout)
│   ├── types/                # AppState, UserInfo, DeviceInfo, QueryRequestPayload, etc.
│   └── utils/                # PrintProgressBar
├── config.toml               # Safe template (CHANGE_ME placeholders), tracked
├── config.test.toml          # Integration test config (safe template)
├── config.{tenant}.toml      # Real credentials, all gitignored via config.*.toml
├── .justfile                 # Windows PowerShell justfile
├── .windows.justfile         # Source for .justfile on Windows
└── .linux.justfile           # Source for .justfile on Linux/macOS
```

## Key rules

### `main()` function — STRICT
**Forbidden inside `main()`:** `for` loops, `csv.NewWriter`, `json.Marshal`, anonymous functions, `os.Create`. Helper naming prefixes:

| Prefix | Purpose | Example |
|--------|---------|---------|
| `fetch*` | API data retrieval with progress | `fetchAllEnrollmentRequests` |
| `write*` | CSV report writer | `writeOneAppToManyUsersCSV` |
| `compute*` | Pure filter/transform | `computeIncludedRequests` |
| `report*` | Report CSV with logging | `reportIncludedRequests` |
| `build*` | Map/set/index construction | `buildExcludeUserIDs` |
| `marshal*` | JSON/serialization | `marshalAppsJSON` |

### Imports — always alias endpoint packages
```go
get_enrollment_requests "msuite-toolkit/pkg/endpoints/get-enrollment-requests"
```
Never use bare imports for `pkg/endpoints` subpackages. Hyphenated directory names → snake_case package names (e.g. `find-user-by-email` → `finduserbyemail`).

### API endpoint pattern (3 layers per package)
1. `GetX(as, payload)` — single page
2. `GetAllX(as, basePayload, progressChan)` — paginated via `pond` pool
3. `GetXWithProgress(as, basePayload)` — progress bar wrapper

URL: `https://{admin_portal_address}/{path}?ctx.user_id={admin_user_id}&request_payload={json}`

Auth: `Authorization: Bearer {token}`

### Config loading
`app.Init("default.csv")` is the first call in every `main()`. Parses `-config` (default `./config.toml`), `-output`, `-h`/`-help`. Sets global `app.AppState`. Access via `as := &app.AppState`.

### CSV writing — pipe-delimited
```go
csvFile, err := os.Create(fileName)
defer func() { csvFile.Close() }()
w := csv.NewWriter(csvFile)
w.Comma = '|'
defer w.Flush()
w.Write(header)
for ... { w.Write(row) }
w.Error()  // mandatory final check
```

### Server-side filter format
```go
[]any{map[string]any{
    "key":      "field_name",
    "operator": "equal_to",
    "value":    "field_value",
}}
```

### Worker pool
Use `pond.NewPool(as.WorkerCount)` for pagination and parallel lookups. `pool.SubmitErr(func() error { ... })` for error-returning tasks. `pool.StopAndWait()` before collecting errors.

## Developer commands (via `just`)
| Command | What it does |
|---------|-------------|
| `just build NAME` | Build `./cmd/<NAME>` → `./dist/<NAME>/<NAME>.exe` + docs + config |
| `just run NAME CONFIG` | `go run ./cmd/<NAME> -config <CONFIG>` |
| `just run NAME CONFIG OUTPUT` | Same with custom `-output` |
| `just test PACKAGE` | `go test -v ./pkg/<PACKAGE>` |
| `just check` | `go build ./cmd/...` |
| `just format` | `go fmt ./...` |

## Testing
- Unit tests for pure functions (`computeIncludedRequests`, `selectInactiveUsers`, `filterAppsByDestination`)
- Integration tests load `../../../config.test.toml` relative to `pkg/endpoints/<package>/` and `t.Skip` when required fields are empty
- Integration tests call real API endpoints — skip by default in CI unless config is populated

## Environment
- Go 1.25.6, only 2 external deps: `github.com/BurntSushi/toml`, `github.com/alitto/pond/v2`
- HTTP client: `InsecureSkipVerify: true`, 30s timeout, 200 max idle conns
- Dev container available in `.devcontainer/` (Debian trixie-slim, Go 1.26.0)
- `*.csv` ignored via `.gitignore`
