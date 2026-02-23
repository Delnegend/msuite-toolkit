# M-Suite Batch Account Inactivate CLI Usage

**Quick steps**

- Ensure M-Suite is open, turned on, and the Admin Portal app is present in the app list.
- Fill `config.toml` (see instructions below).
- Prepare the input file: create a file named `users.txt` with one user ID per line (no header). Example:

```
12345
67890
```

- Open a terminal in this directory (right-click this folder and select "Open in Terminal") and run the tool:

```
./inactive-users.exe -input users.txt
```

**Flags**

- `-config`: path to config file (default: `./config.toml`)
- `-output`: path to output CSV file (default: run `inactive-users.exe -help` to see default name). The output file is the result log: a pipe-separated CSV with header `UserID|Result` where `Result` is `OK` or an error message.
- `-input`: path to input file (required). Each line of the input file should contain a single user ID to deactivate.
- `-h` or `-help`: show help

**config.toml - how to fill**

- `bearer_token`: open the Admin Portal in your browser, open Developer Tools -> Application (or Storage) -> Local Storage -> select the admin portal origin -> find the key `admin_portal_access_token` and copy its value.
- `admin_user_id`: in the Admin Portal navigate to Identity > Users, Groups & Unit > Users, find the currently-logged-in admin user, click the result, then copy the `User ID` from Basic info.
- `admin_portal_address`: set to the Admin Portal host:port you are using (for example `10.0.0.1:9443`).

**Run notes**

- After filling `config.toml`, run `./inactive-users.exe -input users.txt`. The tool will call the Admin Portal to lock each user listed in `users.txt` and write results to the output CSV (see `-output`).
- The output CSV contains rows like:

```
UserID|Result
12345|OK
67890|request failed: unexpected status code: 500
```
