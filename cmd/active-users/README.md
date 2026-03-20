# M-Suite Active Users CLI

## Description
This tool fetches all users from the Admin Portal and writes a JSON file containing users who are considered "active".

Active filtering options (configured in `config.toml`):
- `last_login_threshold_in_month`: discard users who have not logged in within the last X months (0 = disabled)
- `organizational_unit_id`: optional OU id to restrict results to a specific organizational unit (empty = ignored)

## Quick steps
- `cmd/active-users/config.tom` is included in this repo and will be used by default.
- Run the tool from this directory:
```
./active-users.exe
```

## Flags
- `-config`: path to config file (default: `./config.toml`)
- `-output`: path to output JSON file (default: `active_users.json`)
- `-h` or `-help`: show help

## config.toml - how to fill
- `bearer_token`: copy from Admin Portal local storage key `admin_portal_access_token`.
- `admin_user_id`: admin user's ID from Basic info in Admin Portal.
- `admin_portal_address`: host:port of the Admin Portal (e.g. `10.0.0.1:9443`).
- `last_login_threshold_in_month`: integer months threshold (0 disables).
- `organizational_unit_id`: OU id to filter users (leave blank to ignore).

## Example
The repository includes `cmd/active-users/config.tom`. `just build` will bundle this
command-specific config into the `dist/` artifact; otherwise the root `config.toml`
is used. To run locally:

```
go run ./cmd/active-users -config ./cmd/active-users/config.toml -output actives.json
```
