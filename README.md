# defenux

Linux security and hardening CLI.

## Status

Early development. Current version: 0.0.2.

## Commands

    defenux port list
    defenux port status 8000
    defenux network interfaces
    defenux network routes

## Development

Build:

    go build -o defenux .

Test:

    go test ./...

The 0.1.0 release will cover port/network inspection, firewall control,
SSH and service hardening, kernel/user/filesystem auditing, and unified
audit reporting.
