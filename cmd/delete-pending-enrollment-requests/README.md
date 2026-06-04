# Pending Enrollment Requests Delete CLI

## Description
This tool fetches all enrollment requests from the Admin Portal, keeps only requests with `pending` status, and deletes them using the bulk delete endpoint.

## Quick steps
- Ensure M-Suite is open, turned on, and the Admin Portal app is present in the app list.
- Fill `config.toml` (see instructions below).

## Run notes

After filling `config.toml`, run:

```
./delete-pending-enrollment-requests.exe
```

Use `-config` to point to a different config file.

## Flags
- `-config`: path to config file (default: `./config.toml`)
- `-h` or `-help`: show help

## config.toml - how to fill
- `bearer_token`: copy from Admin Portal local storage key `admin_portal_access_token`.
- `admin_user_id`: admin user's ID from Basic info in Admin Portal.
- `admin_portal_address`: host:port of the Admin Portal (for example `10.0.0.1:9443`).
- `worker_count`: optional concurrency setting for paginated requests (default: `100`).
- `dry_run`: set to `true` to preview deletions without executing (writes `to-be-deleted-pending-enrollment-requests.csv`); set to `false` to actually delete (writes `deleted-pending-enrollment-requests.csv`).

Reports
- `deleted-pending-enrollment-requests.csv` / `to-be-deleted-pending-enrollment-requests.csv` — pending requests that were (or would be) deleted.
- `excluded-pending-enrollment-requests.csv` — requests that were kept because their status is not `Pending`.

## Example
Use the `config.toml` in the project root. To run locally:

```
./delete-pending-enrollment-requests.exe -config ./config.toml
```