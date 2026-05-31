# M-Suite Active Users CLI

## Description
This tool fetches all users from the Admin Portal and writes a JSON file containing users who are considered "active".

Active filtering options (configured in `config.toml`):
- `last_login_threshold_in_month`: discard users who have not logged in within the last X months (0 = disabled)
- `organizational_unit_id`: optional OU id to restrict results to a specific organizational unit (empty = ignored)

## Quick steps
- A `config.toml` is included in this folder and is used by default.

## Run notes

After filling `config.toml`, run:

```
./get-active-users.exe
```

Use `-config` to point to a different config file and `-output` to change the output file.

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
Use the `config.toml` in this folder. To run locally:

```
./get-active-users.exe -config ./config.toml -output actives.json
```
