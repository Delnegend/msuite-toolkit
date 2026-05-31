# M-Suite Batch Account Inactivate CLI Usage

## Description
This tool fetches all users from the Admin Portal, selects the accounts whose `last_login_time` is older than the configured threshold, and inactivates those accounts. It generates a CSV report indicating the success or failure of each inactivation.

## Quick steps

- Ensure M-Suite is open, turned on, and the Admin Portal app is present in the app list.
- A `config.toml` is included in this folder and is used by default.
- Run the executable:

```
./set-inactive-users.exe
```

## Flags

- `-config`: path to config file (default: `./config.toml`)
- `-output`: path to output CSV file (default: `inactive_users.csv`). The output file is the result log: a pipe-separated CSV with header `UserID|Result` where `Result` is `OK` or an error message.
- `-h` or `-help`: show help

## config.toml - how to fill

- `bearer_token`: open the Admin Portal in your browser, open Developer Tools -> Application (or Storage) -> Local Storage -> select the admin portal origin -> find the key `admin_portal_access_token` and copy its value.
- `admin_user_id`: in the Admin Portal navigate to Identity > Users, Groups & Unit > Users, find the currently-logged-in admin user, click the result, then copy the `User ID` from Basic info.
- `admin_portal_address`: set to the Admin Portal host:port you are using (for example `10.0.0.1:9443`).
- `last_login_threshold_in_month`: users whose `last_login_time` is older than this many months are considered inactive.
- `dry_run`: when `true`, the tool only lists the users that would be inactivated and does not lock any account.
- `include_users_with_unknown_last_login`: when `true`, users without a readable `last_login_time` are included in the would-be-inactive list. The default is `false`.

## Guardrail

If `last_login_threshold_in_month` is less than `3` and `dry_run` is `false`, the tool requires two confirmations in order:

1. Confirm the dangerous threshold.
2. Confirm that you have read the would-be-inactive users list.

The exact confirmation texts are:

```
I WANT TO INACTIVE USERS WITH LAST LOGIN OLDER THAN X MONTHS
I HAVE READ THE WOULD-BE-INACTIVE USERS LIST AND WANT TO PROCEED
```

`X` is replaced with the configured threshold value.

## Run notes

After filling `config.toml`, run:

```
./set-inactive-users.exe
```

Use `-config` to point to a different config file and `-output` to change the output file. The tool fetches all users, selects the accounts older than the configured threshold, and locks them unless `dry_run` is enabled. When `dry_run` is `false`, the dangerous-threshold prompt runs first if the threshold is below 3 months, then the list-read confirmation runs immediately after. The output CSV contains rows like:

```
UserID|Result
12345|OK
67890|request failed: unexpected status code: 500
```
