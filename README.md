# **gdenv** ![GitHub release (with filter)](https://img.shields.io/github/v/release/coffeebeats/gdenv) ![GitHub](https://img.shields.io/github/license/coffeebeats/gdenv) [![Build Status](https://img.shields.io/github/actions/workflow/status/coffeebeats/gdenv/check-commit.yml?branch=main)](https://github.com/coffeebeats/gdenv/actions?query=branch%3Amain+workflow%3Acheck) [![codecov](https://codecov.io/gh/coffeebeats/gdenv/graph/badge.svg)](https://codecov.io/gh/coffeebeats/gdenv)

A single-purpose, CI-friendly command-line interface for managing Godot versions. Inspired by [pyenv](https://github.com/pyenv/pyenv), [rbenv](https://github.com/rbenv/rbenv), and [volta](https://github.com/volta-cli/volta).

## **Getting started**

These instructions will help you install `gdenv` and pin projects (or your system) to specific versions of _Godot_.

### **Example usage**

After following the [installation instructions](#installation), the following are example usages of `gdenv`:

#### Install a global (system-wide) _Godot_ version

```sh
gdenv pin -ig 4.0
```

#### Pin a project to a specific _Godot_ version

```sh
# Omit the `--path` option to pin the current directory. The `-i` flag instructs `gdenv` to download the pinned version to its cache.
gdenv pin -i --path /path/to/project 4.0
```

#### Vendor the _Godot_ source code

```sh
# Omit the `--path` option to vendor to `./godot-4.0-stable`.
gdenv vendor --path /path/to/source 4.0
```

### **Installation**

The easiest way to install `gdenv` is by using the pre-built binaries. These can be manually downloaded and configured, but automated installation scripts are provided and recommended.

> ⚠️ **WARNING:** It's good practice to inspect an installation script prior to execution. The scripts are included in this repository and can be reviewed prior to use.

#### **Linux/MacOS — `sh` (recommended)**

> ❕ **NOTE:** If you're using [Git BASH for Windows](https://gitforwindows.org/), use these instructions instead of [Windows (powershell)](#windows--powershell-recommended).

```sh
curl https://raw.githubusercontent.com/coffeebeats/gdenv/main/scripts/install.sh | sh
```

#### **Windows — `powershell` (recommended)**

> ❕ **NOTE:** In order to run scripts in PowerShell, the [execution policy](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_execution_policies) must _not_ be `Restricted`. Consider running the following command
> if you encounter `UnauthorizedAccess` errors when following these instructions. See [Set-ExecutionPolicy](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.security/set-executionpolicy) documentation for details.
>
> ```sh
> Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope LocalMachine
> ```

```sh
Invoke-WebRequest `
    -UseBasicParsing `
    -Uri "https://raw.githubusercontent.com/coffeebeats/gdenv/main/scripts/install.ps1" `
    -OutFile "./install-gdenv.ps1"; `
    &"./scripts/install-gdenv.ps1"
```

#### **Manual (not recommended)**

> ❕ **NOTE:** The instructions below provide `bash`-specific commands for a _Linux_-based system. While these won't work in _PowerShell_, the process will be similar.

1. Download a prebuilt binary from the corresponding GitHub release.

    ```sh
    # Set '$VERSION', '$OS', and '$ARCH' to the desired values.
    VERSION=0.0.0 OS=linux ARCH=x86_64; \
    curl -LO https://github.com/coffeebeats/gdenv/releases/download/v$VERSION/gdenv-$VERSION-$OS-$ARCH.tar.gz
    ```

2. Extract the downloaded archive.

    ```sh
    # Set '$GDENV_HOME' to the desired location.
    GDENV_HOME=$HOME/.gdenv; \
    mkdir -p $GDENV_HOME/bin && \
    tar -C $GDENV_HOME/bin -xf gdenv-$VERSION-$OS-$ARCH.tar.gz
    ```

3. Export the `$GDENV_HOME` variable and add `$GDENV_HOME/bin` to `$PATH`.

    ```sh
    # In '.bashrc', or something similar ('$GDENV_HOME' can be customized).
    export GDENV_HOME="$HOME/.gdenv"
    export PATH="$GDENV_HOME/bin:$PATH"
    ```

#### **Compile from source (not recommended)**

`gdenv` is a Go project and can be installed using `go install`. This option is not recommended as it requires having the Go toolchain installed, it's slower than downloading a prebuilt binary, and there may be instability due to using a different version of Go than it was developed with.

> NOTE: You will need to somehow set the installed `gdenv-shim` binary as your system's `godot` command (consider using a symbolic link). This is done automatically by the recommended installation methods listed above.

```sh
go install github.com/coffeebeats/gdenv/cmd/gdenv@latest
```

## **Documentation**

### **How it works**

The `gdenv` application maintains a cache of downloaded _Godot_ executables (typically `$HOME/.gdenv`) and provides a shim which should be set to the system's `godot` executable. The shim will examine the current directory from which the `godot` command is invoked (handling a `--path` option as well) and delegate to the correct version of _Godot_.

In order to track pinned versions of _Godot_, the `pin` subcommand will place a `.godot-version` file in the specified directory (or within `$GDENV_HOME` if pinning a global version with `-g`). This is what the `godot` shim will use to determine the correct _Godot_ version.

### **gdenv `pin`**

Sets the _Godot_ version globally or for a specific directory.

**Options:**

- **`-g`**, **`--global`** — pin the system version (cannot be used with `-p`)
- **`-i`**, **`--install`** — installs the specified version of _Godot_ if missing
- **`-f`**, **`--force`** — forcibly overwrite an existing cache entry (only used with `-i`)
- **`-p`**, **`--path <PATH>`** — pin the specified path (cannot be used with `-g`)
  - Default value: `$PWD` (current working directory)

**Arguments:**

- **`<VERSION>`** — the specific version string to install (must be exact)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

### **gdenv `unpin`**

Removes a `Godot` version pin from the system or specified directory.

**Options:**

- **`-g`**, **`--global`** — unpin the system version (cannot be used with `-p`)
- **`-p`**, **`--path <PATH>`** — unpin the specified path (cannot be used with `-g`)
  - Default value: `$PWD` (current working directory)

### **gdenv `install`**

Downloads and caches a specific version of _Godot_.

**Options:**

- **`-f`**, **`--force`** — forcibly overwrite an existing cache entry
- **`-g`**, **`--global`** — pin the system version (cannot be used with `-p`)
- **`-p`**, **`--path <PATH>`** — determine the version from the pinned `PATH` (ignores the global pin)
- **`-s`**, **`--src`**, **`--source`** — install source code instead of an executable (cannot be used with `-g`)

**Arguments:**

- **`[VERSION]`** — the specific version string to install (must be exact)
  - Default value: Resolves the pinned version at `$PWD` (ignoring the global pin)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

### **gdenv `vendor`**

Download the _Godot_ source code to the specified directory.

**Options:**

- **`-f`**, **`--force`** — forcibly overwrite an existing cache entry
- **`-o`**, **`--out`** — download the source code into `OUT` (will overwrite conflicting files)
  - Default value: `$PWD/./godot-<VERSION>`
- **`-p`**, **`--path <PATH>`** — determine the version from the pinned `PATH` (ignores the global pin)
  - Default value: `$PWD` (current working directory)

**Arguments:**

- **`[VERSION]`** — the specific version string to install (must be exact and cannot be used with `-p`)
  - Default value: Resolves the pinned version at `$PWD` (ignoring the global pin)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

### **gdenv `uninstall`**

Removes the specified version of _Godot_ from the `gdenv` download cache.

**Options:**

- **`-a`**, **`--all`** — uninstall all versions of _Godot_ (ignores source code without `-s`)
- **`-s`**, **`--src`**, **`--source`** — uninstall source code versions

**Arguments:**

- **`[VERSION]`** — the specific version string to install (must be exact; omit if using `-a`)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

### **gdenv `ls`/`list`**

Prints the path and version of all of the installed versions of _Godot_.

**Options:**

- **`-a`**, **`--all`** — list executable _and_ source code versions
- **`-s`**, **`--src`**, **`--source`** — list source code versions

### **gdenv `which`**

Prints the path to the _Godot_ executable which would be used in the specified directory.

**Options:**

- **`-p`**, **`--path <PATH>`** — the specified path to check
  - Default value: `$PWD` (current working directory)

### **gdenv `completions`**

> ⚠️ **WARNING:** This command is not yet implemented.

Provides shell completions for the `gdenv` CLI application.

**Options:**

- **`-o`**, **`--output <OUT_FILE>`** — File to write the completions to
  - Default value: `stdout`

**Arguments:**

- **`<SHELL>`** — the specific version string to install (must be exact)
  - Supported values:
    - `bash`
    - `fish`
    - `powershell`
    - `zsh`

## **Development**

The following instructions outline how to get the project set up for local development:

1. [Follow the instructions](https://go.dev/doc/install) to install Go (see [go.mod](./go.mod) for the minimum required version).
2. Clone the [coffeebeats/gdenv](https://github.com/coffeebeats/gdenv) repository.
3. Install the [required tools](./tools.go) using the following command:

    ```sh
    cat tools.go | grep _ | grep -v '//' | awk -F'"' '{print $2}' | xargs -tI % go install %
    ```

When submitting code for review, ensure the following requirements are met:

1. The project is correctly formatted using [go fmt](https://go.dev/blog/gofmt):

    ```sh
    go fmt ./...
    ```

2. All [golangci-lint](https://golangci-lint.run/) linter warnings are addressed:

    ```sh
    go fmt ./...
    ```

3. All unit tests pass and no data races are found:

    ```sh
    go test -race ./...
    ```

4. The `gdenv` and `gdenv-shim` binaries successfully compile (release artifacts will be available at `./dist`):

    ```sh
    goreleaser release --clean --skip-publish --snapshot
    ```

## **Contributing**

All contributions are welcome! Feel free to file [bugs](https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=bug&projects=&template=%F0%9F%90%9B-bug-report.md&title=) and [feature requests](https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=enhancement&projects=&template=%F0%9F%99%8B-feature-request.md&title=) and/or open pull requests.

## **Version history**

See [CHANGELOG.md](https://github.com/coffeebeats/gdenv/blob/main/CHANGELOG.md).

## **License**

[MIT License](https://github.com/coffeebeats/gdenv/blob/main/LICENSE)
