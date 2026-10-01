# defenux

Linux security and hardening CLI.

## Install

Latest Linux release:

    curl -fsSL https://raw.githubusercontent.com/zhravan/defenux/main/install.sh | sh

Install a specific version:

    DEFENUX_VERSION=0.0.2 curl -fsSL https://raw.githubusercontent.com/zhravan/defenux/main/install.sh | sh

The installer detects the Linux CPU architecture, downloads the matching static
binary from GitHub Releases, verifies its SHA-256 checksum, and installs it to
/usr/local/bin or ~/.local/bin.

Supported Linux architectures:

- amd64
- arm64
- armv6
- armv7
- 386
- ppc64le
- s390x
- riscv64
- loong64

## Commands

    defenux port list
    defenux port status 8000
    defenux network interfaces
    defenux network routes
    defenux firewall status
    defenux firewall list
    defenux firewall allow 22/tcp
    defenux firewall deny 23/tcp
    defenux ssh status
    defenux ssh config

## Development

Build:

    go build -o defenux .

Test:

    go test ./...

Source builds report dev. Release builds embed the Git tag version.

The 0.1.0 release will cover port/network inspection, firewall control,
SSH and service hardening, kernel/user/filesystem auditing, and unified
audit reporting.
