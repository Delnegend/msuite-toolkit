# M-Suite User Devices Data Extraction CLI Usage

## Description
This tool extracts a list of all devices associated with each user in the system. The output CSV includes User ID, Email, Device ID, Device Name, Device Type, and the Last Used timestamp for each device.

## Quick steps
- Ensure M-Suite is open, turned on, and the Admin Portal app is present in the app list.
- Fill `config.toml` (see instructions below).
- Open a terminal in this directory (right-click this folder and select "Open in Terminal") and run the tool:
```
./get-user-devices.exe
```

## Flags
- `-config`: path to config file (default: `./config.toml`)
- `-output`: path to output CSV file (default: `user_devices.csv`)
- `-h` or `-help`: show help

## config.toml - how to fill
- `bearer_token`: open the Admin Portal in your browser, open Developer Tools -> Application (or Storage) -> Local Storage -> select the admin portal origin -> find the key `admin_portal_access_token` and copy its value.
- `admin_user_id`: in the Admin Portal navigate to Identity > Users, Groups & Unit > Users, find the currently-logged-in admin user, click the result, then copy the `User ID` from Basic info.
- `admin_portal_address`: set to the Admin Portal host:port you are using (for example `10.0.0.1:9443`).
- `organizational_unit_id`: optional OU id to restrict user results to a specific organizational unit (leave blank to ignore).

## Run notes

After filling `config.toml`, run:

```
./get-user-devices.exe
```

Use `-config` to point to a different config file and `-output` to change the output file. The output CSV is pipe-separated and includes the following columns: `UserID`, `UserEmail`, `DeviceID`, `DeviceName`, `DeviceType`, `LastUsed`.

