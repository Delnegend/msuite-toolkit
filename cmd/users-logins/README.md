# M-Suite Users MFA & Logins Data Extraction CLI Usage

## Description
This tool extracts MFA status and a history of failed login attempts for each user. The output CSV includes User ID, Email, MFA configuration (JSON), and Failed Logins (JSON).

## Quick steps
- Ensure M-Suite is open, turned on, and the Admin Portal app is present in the app list.
- Fill `config.toml` (see instructions below).
- Open a terminal in this directory (right-click this folder and select "Open in Terminal") and run the tool:
```
./users-logins.exe
```

## Flags
- `-config`: path to config file (default: `./config.toml`)
- `-output`: path to output CSV file (default: `users_logins.csv`)
- `-h` or `-help`: show help

## config.toml - how to fill
- `bearer_token`: open the Admin Portal in your browser, open Developer Tools -> Application (or Storage) -> Local Storage -> select the admin portal origin -> find the key `admin_portal_access_token` and copy its value.
- `admin_user_id`: in the Admin Portal navigate to Identity > Users, Groups & Unit > Users, find the currently-logged-in admin user, click the result, then copy the `User ID` from Basic info.
- `admin_portal_address`: set to the Admin Portal host:port you are using (for example `10.0.0.1:9443`).

## Run notes
- After filling `config.toml`, run `./users-logins.exe`. The default output file will be created next to the tool.
- Use `-c` to point to a different config file and `-o` to choose another output file name.

## Example output
```
User ID|Email|MFA|Failed logins
12345|user@example.com|{"type": "totp", "enabled": true}|[{"timestamp": "2023-10-01T10:00:00Z", "ip": "1.2.3.4"}]
```
