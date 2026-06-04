# Enrollment Requests Delete CLI

## Description
This tool fetches all enrollment requests for a given organizational unit (OU) from the Admin Portal and deletes them using the bulk delete endpoint. Unlike `delete-pending-enrollment-requests`, this tool does **not** filter by status — it deletes all requests (Pending, Approved, Rejected, etc.) found for the OU.

> `organizational_unit_id` is **required** in `config.toml`.

## Quick steps
- Ensure M-Suite is open, turned on, and the Admin Portal app is present in the app list.
- Fill `config.toml` in the root directory (see instructions below).

## Run notes

After filling `config.toml`, run:

```
./delete-enrollment-requests.exe
```

Use `-config` to point to a different config file.

## Flags
- `-config`: path to config file (default: `./config.toml`)
- `-h` or `-help`: show help

## config.toml - how to fill
- `bearer_token`: copy from Admin Portal local storage key `admin_portal_access_token`.
- `admin_user_id`: admin user's ID from Basic info in Admin Portal.
- `admin_portal_address`: host:port of the Admin Portal (for example `10.0.0.1:9443`).
- `organizational_unit_id`: **required** — the OU whose enrollment requests should be deleted.
- `dry_run`: set to `true` to preview deletions without executing (writes `to-be-deleted-enrollment-requests.csv`); set to `false` to actually delete (writes `deleted-enrollment-requests.csv`).
- `exclude_emails`: optional list of user emails whose enrollment requests should be kept (not deleted). For example: `exclude_emails = ["admin@example.com"]`.
- `worker_count`: optional concurrency setting for paginated requests (default: `100`).

Reports
- `deleted-enrollment-requests.csv` / `to-be-deleted-enrollment-requests.csv` — requests that were (or would be) deleted.
- `excluded-enrollment-requests.csv` — requests that were kept because the user's email is in `exclude_emails` or because the device type is not a desktop (Windows, macOS, Linux).

## Example
To run from the project root:

```
./delete-enrollment-requests.exe -config ./config.toml
```