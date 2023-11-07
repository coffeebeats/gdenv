# **gdenv** ![GitHub release (with filter)](https://img.shields.io/github/v/release/coffeebeats/gdenv) ![GitHub](https://img.shields.io/github/license/coffeebeats/gdenv) [![Build Status](https://img.shields.io/github/actions/workflow/status/coffeebeats/gdenv/check-commit.yml?branch=main)](https://github.com/coffeebeats/gdenv/actions?query=branch%3Amain+workflow%3Acheck) [![codecov](https://codecov.io/gh/coffeebeats/gdenv/graph/badge.svg)](https://codecov.io/gh/coffeebeats/gdenv)

A single-purpose, CI-friendly command-line interface for managing Godot editor versions. Inspired by [pyenv](https://github.com/pyenv/pyenv), [rbenv](https://github.com/rbenv/rbenv), and [volta](https://github.com/volta-cli/volta).

## **How it works**

The `gdenv` application maintains a cache of downloaded _Godot_ executables (typically `$HOME/.gdenv`) and provides a [shim](https://en.wikipedia.org/wiki/Shim_(computing)) which should be set to the system's `godot` executable. This shim intercepts normal use of the `godot` command and, based on the directory from which the `godot` command is invoked (or the `--path` option), transparently invokes the correct version of _Godot_ with the provided arguments.

In order to track pinned versions of _Godot_, the `pin` subcommand will place a `.godot-version` file in the specified directory (or within `$GDENV_HOME` if pinning a global version with `-g`). This is what the `godot` shim will use to determine the correct _Godot_ version.

## **Getting started**

These instructions will help you install `gdenv` and pin projects (or your system) to specific versions of _Godot_.

### **Example usage**

> NOTE: For _Mono_-flavored builds, see [Version selection (C#/_Mono_ support)](#version-selection-cmono-support).

#### **Install a global (system-wide) _Godot_ version**

```sh
gdenv pin -ig 4.0
```

#### **Pin a project to a specific _Godot_ version**

```sh
# Omit the `--path` option to pin the current directory;
#   the `-i` flag instructs `gdenv` to download the pinned version to its cache.
gdenv pin -i --path /path/to/project 4.0
```

### **Installation**

See [docs/installation.md](./docs/installation.md#installation) for detailed instructions on how to download `gdenv`.

## **API Reference**

### **Commands**

See [docs/commands.md](./docs/commands.md) for a detailed reference on how to use each command.

#### Manage installed versions

- **[install](./docs/commands.md#gdenv-install)** — `gdenv install [OPTIONS] [VERSION]`
- **[uninstall](./docs/commands.md#gdenv-uninstall)** — `gdenv uninstall [OPTIONS] [VERSION]`
- **[vendor](./docs/commands.md#gdenv-vendor)** — `gdenv vendor [OPTIONS] [VERSION]`

#### Pin projects/set system default

- **[pin](./docs/commands.md#gdenv-pin)** — `gdenv pin [OPTIONS] <VERSION>`
- **[unpin](./docs/commands.md#gdenv-unpin)** — `gdenv unpin [OPTIONS]`

#### Inspect versions

- **[ls/list](./docs/commands.md#gdenv-lslist)** — `gdenv ls [OPTIONS]`
- **[which](./docs/commands.md#gdenv-which)** — `gdenv which [OPTIONS]`

### **Platform selection**

By default `gdenv` will install _Godot_ executables for the host platform (i.e. the system `gdenv` is running on). To change which platform `gdenv` selections, the following environment variables can be set in front of any `gdenv` command:

> ❕ **NOTE:** These options are meant to circumvent incorrect platform detection by `gdenv` or facilitate installing different _Godot_ editor versions in a CI environment. Most users will not need to set these when using `gdenv` locally.

- `GDENV_OS` - set the target operating system (still uses the host's CPU architecture)
- `GDENV_ARCH` - set the target CPU architecture (still uses the host's operating system)
- `GDENV_PLATFORM` - set the literal string suffix of the _Godot_ editor (e.g. `macos.universal` or `win64`)

### **Version selection (C#/_Mono_ support)**

`gdenv` considers _Mono_ variants of _Godot_ to be part of the version and not the platform. As such, to have `gdenv` install Mono builds of _Godot_ editors all version specifications should be suffixed with `stable_mono` (e.g. `gdenv pin 4.0-stable_mono` or `gdenv install 4.1.1-stable_mono`). Although `gdenv` normally assumes a `stable` release if the label is omitted, _Mono_ builds must be explicitly specified.

However, to simplify use of `gdenv` when _Mono_ builds are desired, the following environment variable can be set to have `gdenv` default to using _Mono_ builds _when the version label is omitted_. A non-_Mono_ build can then be specified by passing a version label of `stable` without the `_mono` suffix.

- `GDENV_DEFAULT_MONO` - set to `1` to have `gdenv` interpret missing version labels as `stable_mono` instead of `stable`

## **Development**

The following instructions outline how to get the project set up for local development:

1. [Follow the instructions](https://go.dev/doc/install) to install Go (see [go.mod](./go.mod) for the minimum required version).
2. Clone the [coffeebeats/gdenv](https://github.com/coffeebeats/gdenv) repository.
3. Install the [required tools](./tools.go) using the following command ([see this article](https://www.alexedwards.net/blog/using-go-run-to-manage-tool-dependencies) for more explanation):

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
    golangci-lint run ./...
    ```

3. All unit tests pass and no data races are found:

    ```sh
    go test -race ./...
    ```

4. The `gdenv` and `gdenv-shim` binaries successfully compile (release artifacts will be available at `./dist`):

    ```sh
    goreleaser release --clean --skip=publish --snapshot
    ```

## **Contributing**

All contributions are welcome! Feel free to file [bugs](https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=bug&projects=&template=bug-report.md&title=) and [feature requests](https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=enhancement&projects=&template=feature-request.md&title=) and/or open pull requests.

## **Version history**

See [CHANGELOG.md](https://github.com/coffeebeats/gdenv/blob/main/CHANGELOG.md).

## **License**

[MIT License](https://github.com/coffeebeats/gdenv/blob/main/LICENSE)
