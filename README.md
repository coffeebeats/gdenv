# **gdenv** ![GitHub release (with filter)](https://img.shields.io/github/v/release/coffeebeats/gdenv?style=flat-square) [![Build Status](https://img.shields.io/github/actions/workflow/status/coffeebeats/gdenv/check-commit.yml?branch=main&style=flat-square)](https://github.com/coffeebeats/gdenv/actions?query=branch%3Amain+workflow%3Acheck) ![GitHub](https://img.shields.io/github/license/coffeebeats/gdenv?style=flat-square)

> ⚠️ **WARNING:** This repository is in its early stages and is under active development. A lot of functionality is missing, and there is no guarantee of API stability.

A single-purpose, CI-friendly command-line interface for managing Godot versions. Inspired by [pyenv](https://github.com/pyenv/pyenv), [rbenv](https://github.com/rbenv/rbenv), and [volta](https://github.com/volta-cli/volta).

## **Getting started**

These instructions will help you install `gdenv` and pin projects (or your system) to specific versions of _Godot_.

### **Example usage**

#### Install a global (system-wide) _Godot_ version

```sh
gdenv pin -ig 4.1.1
```

#### Pin a project to a specific _Godot_ version

```sh
gdenv pin --path /path/to/project -i 4.1.1
```

### **Installation**

The easiest way to install `gdenv` is by using the pre-built binaries. These can be manually downloaded and configured, but automated installation scripts are provided and recommended.

> ⚠️ **WARNING:** It's good practice to inspect an installation script prior to execution. The scripts are included in this repository and can be reviewed prior to use.

#### **Linux/MacOS — bash (recommended)**

> ❕ **NOTE:** If you're using [Git BASH for Windows](https://gitforwindows.org/), use these instructions instead of [Windows (powershell)](#windows--powershell-recommended).

```sh
curl https://raw.githubusercontent.com/coffeebeats/gdenv/main/install.sh | bash
```

#### **Windows — powershell (recommended)**

> ❕ **NOTE:** In order to run scripts in PowerShell, the [execution policy](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_execution_policies) must _not_ be `Restricted`. Consider running the following command
> if you encounter `UnauthorizedAccess` errors when following these instructions. See [Set-ExecutionPolicy](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.security/set-executionpolicy) documentation for details.
>
> ```sh
> Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope LocalMachine
> ```

```sh
Invoke-WebRequest `
    -UseBasicParsing `
    -Uri "https://raw.githubusercontent.com/coffeebeats/gdenv/main/install.ps1" `
    -OutFile "./install-gdenv.ps1"; `
    &"./install-gdenv.ps1"
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

TODO: Provide instructions for compiling from source.

## **Documentation**

### **How it works**

The `gdenv` application maintains a cache of downloaded _Godot_ executables (typically `$HOME/.gdenv`) and provides a shim which should be set to the system's `godot` executable. The shim will examine the current directory from which the `godot` command is invoked (handling a `--path` option as well) and delegate to the correct version of _Godot_.

In order to track pinned versions of _Godot_, the `pin` subcommand will place a `.godot-version` file in the specified directory (or within `$GDENV_HOME` if pinning a global version with `-g`). This is what the `godot` shim will use to determine the correct _Godot_ version.

### **gdenv `pin`**

Sets the _Godot_ version globally or for a specific directory.

**Options:**

- **`-i`**, **`--install`** — installs the specified version of _Godot_ if missing
- **`-g`**, **`--global`** — pin the system version (cannot be used with `-p`)
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

**Arguments:**

- **`<VERSION>`** — the specific version string to install (must be exact)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

### **gdenv `uninstall`**

Removes the specified version of _Godot_ from the `gdenv` download cache.

**Options:**

- **`-a`**, **`--all`** — uninstall all versions of _Godot_ in the `gdenv` cache

**Arguments:**

- **`[VERSION]`** — the specific version string to install (must be exact; omit if using `-a`)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

### **gdenv `ls`**

Prints the path and version of all of the installed versions of _Godot_.

### **gdenv `which`**

Prints the path to the _Godot_ executable which would be used in the specified directory.

**Options:**

- **`-p`**, **`--path <PATH>`** — the specified path to check
  - Default value: `$PWD` (current working directory)

### **gdenv `completions`**

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

TODO: Provide development environment setup instructions.

## **Contributing**

All contributions are welcome! Feel free to open pull request or file [bugs](https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=bug&projects=&template=%F0%9F%90%9B-bug-report.md&title=) and [feature requests](https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=enhancement&projects=&template=%F0%9F%99%8B-feature-request.md&title=).

## **Version history**

See [CHANGELOG.md](https://github.com/coffeebeats/gdenv/blob/main/CHANGELOG.md).

## **License**

[MIT License](https://github.com/coffeebeats/gdenv/blob/main/LICENSE)
