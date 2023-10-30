# **gdenv** ![GitHub release (with filter)](https://img.shields.io/github/v/release/coffeebeats/gdenv) ![GitHub](https://img.shields.io/github/license/coffeebeats/gdenv) [![Build Status](https://img.shields.io/github/actions/workflow/status/coffeebeats/gdenv/check-commit.yml?branch=main)](https://github.com/coffeebeats/gdenv/actions?query=branch%3Amain+workflow%3Acheck) [![codecov](https://codecov.io/gh/coffeebeats/gdenv/graph/badge.svg)](https://codecov.io/gh/coffeebeats/gdenv)

A single-purpose, CI-friendly command-line interface for managing Godot versions. Inspired by [pyenv](https://github.com/pyenv/pyenv), [rbenv](https://github.com/rbenv/rbenv), and [volta](https://github.com/volta-cli/volta).

## **Getting started**

These instructions will help you install `gdenv` and pin projects (or your system) to specific versions of _Godot_.

### **Example usage**

> NOTE: For _Mono_-flavored builds, see [Version selection (C#/_Mono_ support)](#version-selection-cmono-support).

#### Install a global (system-wide) _Godot_ version

```sh
gdenv pin -ig 4.0
```

#### Pin a project to a specific _Godot_ version

```sh
# Omit the `--path` option to pin the current directory;
#   the `-i` flag instructs `gdenv` to download the pinned version to its cache.
gdenv pin -i --path /path/to/project 4.0
```

### **Installation**

The easiest way to install `gdenv` is by using the pre-built binaries. These can be manually downloaded and configured, but automated installation scripts are provided and recommended.

See the full [installation instructions](./docs/installation.md) for additional options for installing `gdenv`.

> ⚠️ **WARNING:** It's good practice to inspect an installation script prior to execution. The scripts are included in this repository and can be reviewed prior to use.

#### **Linux/MacOS — `sh`**

```sh
curl https://raw.githubusercontent.com/coffeebeats/gdenv/main/scripts/install.sh | sh
```

#### **Windows - Git BASH for Windows**

> ❕ **NOTE:** If you're using [Git BASH for Windows](https://gitforwindows.org/), use the [Linux/MacOS instructions](#linuxmacos--sh).

#### **Windows — `powershell`**

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

## **Documentation**

### Commands

See [Commands](./docs/commands.md) for more explanation about how to use `gdenv`.

### **How it works**

The `gdenv` application maintains a cache of downloaded _Godot_ executables (typically `$HOME/.gdenv`) and provides a shim which should be set to the system's `godot` executable. The shim will examine the current directory from which the `godot` command is invoked (handling a `--path` option as well) and delegate to the correct version of _Godot_.

In order to track pinned versions of _Godot_, the `pin` subcommand will place a `.godot-version` file in the specified directory (or within `$GDENV_HOME` if pinning a global version with `-g`). This is what the `godot` shim will use to determine the correct _Godot_ version.

### Platform selection

By default `gdenv` will install _Godot_ executables for the host platform (i.e. the system `gdenv` is running on). To change which platform `gdenv` selections, the following environment variables can be set in front of any `gdenv` command:

> ❕ **NOTE:** These options are meant to circumvent incorrect platform detection by `gdenv` or facilitate installing different _Godot_ editor versions in a CI environment. Most users will not need to set these when using `gdenv` locally.

- `GDENV_OS` - set the target operating system (still uses the host's CPU architecture)
- `GDENV_ARCH` - set the target CPU architecture (still uses the host's operating system)
- `GDENV_PLATFORM` - set the literal string suffix of the _Godot_ editor (e.g. `macos.universal` or `win64`)

### Version selection (C#/_Mono_ support)

`gdenv` considers _Mono_ variants of _Godot_ to be part of the version and not the platform. As such, to have `gdenv` install Mono builds of _Godot_ editors all version specifications should be suffixed with `stable_mono` (e.g. `gdenv pin 4.0-stable_mono` or `gdenv install 4.1.1-stable_mono`). Although `gdenv` normally assumes a `stable` release if the label is omitted, _Mono_ builds must be explicitly specified.

However, to simplify use of `gdenv` when _Mono_ builds are desired, the following environment variable can be set to have `gdenv` default to using _Mono_ builds _when the version label is omitted_. A non-_Mono_ build can then be specified by passing a version label of `stable` without the `_mono` suffix.

- `GDENV_MONO_DEFAULT` - set to something truthy (e.g. `1`) to have `gdenv` interpret missing version labels as `stable_mono` instead of `stable`

## **Development**

The following instructions outline how to get the project set up for local development:

1. [Follow the instructions](https://go.dev/doc/install) to install Go (see [go.mod](./go.mod) for the minimum required version).
2. Clone the [coffeebeats/gdenv](https://github.com/coffeebeats/gdenv) repository.
3. Install the [required tools](./tools.go) [using the following command](https://www.alexedwards.net/blog/using-go-run-to-manage-tool-dependencies):

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

All contributions are welcome! Feel free to file [bugs](https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=bug&projects=&template=bug-report.md&title=) and [feature requests](https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=enhancement&projects=&template=feature-request.md&title=) and/or open pull requests.

## **Version history**

See [CHANGELOG.md](https://github.com/coffeebeats/gdenv/blob/main/CHANGELOG.md).

## **License**

[MIT License](https://github.com/coffeebeats/gdenv/blob/main/LICENSE)
