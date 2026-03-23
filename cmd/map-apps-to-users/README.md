# M-Suite Apps to Users Mapping Data Extraction CLI Usage

## Description
This tool extracts information about which users have access to which applications. It generates two CSV files:
- `ONE-to-MANY`: A mapping of each application to a list of users who have access to it.
- `ONE-to-ONE`: A direct mapping of each user to each application they have access to.

## Quick steps
- Ensure M-Suite is open, turned on, and the Admin Portal app is present in the app list.
- Fill `config.toml` (see instructions below).
- Open a terminal in this directory (right-click this folder and select "Open in Terminal") and run the tool:
```
./map-apps-to-users.exe
```

## Flags
- `-config`: path to config file (default: `./config.toml`)
- `-output`: path to output CSV file name suffix (default: `apps_to_users.csv`). Two files will be created: `ONE-to-MANY_<output>` and `ONE-to-ONE_<output>`.
- `-h` or `-help`: show help

## config.toml - how to fill
- `bearer_token`: open the Admin Portal in your browser, open Developer Tools -> Application (or Storage) -> Local Storage -> select the admin portal origin -> find the key `admin_portal_access_token` and copy its value.
- `admin_user_id`: in the Admin Portal navigate to Identity > Users, Groups & Unit > Users, find the currently-logged-in admin user, click the result, then copy the `User ID` from Basic info.
- `admin_portal_address`: set to the Admin Portal host:port you are using (for example `10.0.0.1:9443`).

## Run notes
- After filling `config.toml`, run `./map-apps-to-users.exe`. The default output files will be created next to the tool.
- The tool generates two files:
  - `ONE-to-MANY_apps_to_users.csv`: Maps each app to a comma-separated list of user IDs.
  - `ONE-to-ONE_apps_to_users.csv`: Maps each app to a single user ID per row.

Use `-c` to point to a different config file and `-o` to choose another output file name suffix.
