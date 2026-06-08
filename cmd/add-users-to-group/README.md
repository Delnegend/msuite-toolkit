# Add Users To Group CLI

## Description
This tool resolves a list of user emails to user IDs and adds those users to a
target group via the bulk add endpoint. Users are added in **batches of 10**,
sent concurrently, with a progress bar.

> Both `group_id` and `emails` are **required** in `config.toml`.

## Quick steps
- Ensure M-Suite is open, turned on, and the Admin Portal app is present in the app list.
- Fill `config.toml` in the root directory (see instructions below).

## Run notes

After filling `config.toml`, run:

```
./add-users-to-group.exe
```

Use `-config` to point to a different config file.

## Flags
- `-config`: path to config file (default: `./config.toml`)
- `-h` or `-help`: show help

## config.toml - how to fill
- `bearer_token`: copy from Admin Portal local storage key `admin_portal_access_token`.
- `admin_user_id`: admin user's ID from Basic info in Admin Portal.
- `admin_portal_address`: host:port of the Admin Portal (for example `10.0.0.1:9443`).
- `group_id`: **required** — the group the users will be added to.
- `emails`: **required** — list of user emails to add. Each email is resolved to a user ID before being added. For example: `emails = ["a@example.com", "b@example.com"]`.
- `dry_run`: set to `true` to preview (resolves emails and writes `to-be-added-users.csv` without adding); set to `false` to actually add (writes `added-users.csv`).
- `worker_count`: optional concurrency setting for email resolution and add batches (default: `100`).

Reports
- `added-users.csv` / `to-be-added-users.csv` — users that were (or would be) added, with their resolved user IDs.
- `unresolved-emails.csv` — emails that could not be matched to a user (written only when some fail to resolve).

## Example
To run from the project root:

```
./add-users-to-group.exe -config ./config.toml
```
