# M-Suite Provision Policies Extraction CLI Usage

## Description
This tool extracts all provision policies from the M-Suite system. The output CSV includes Policy ID, Policy Name, Created Time, and the full policy configuration in JSON format.

## Quick steps
- Ensure M-Suite is open, turned on, and the Admin Portal app is present in the app list.
- Fill `config.toml` (see instructions below).
- Open a terminal in this directory (right-click this folder and select "Open in Terminal") and run the tool:
```
./get-provision-policies.exe
```

## Flags
- `-config`: path to config file (default: `./config.toml`)
- `-output`: path to output CSV file (default: `provision_policies.csv`)
- `-h` or `-help`: show help

## Run notes

After filling `config.toml`, run:

```
./get-provision-policies.exe
```

Use `-config` to point to a different config file and `-output` to change the output file.
