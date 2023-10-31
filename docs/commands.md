# Commands

**Supported commands:**

- **[gdenv completions](#gdenv-completions)**
- **[gdenv install](#gdenv-install)**
- **[gdenv ls/list](#gdenv-lslist)**
- **[gdenv pin](#gdenv-pin)**
- **[gdenv uninstall](#gdenv-uninstall)**
- **[gdenv unpin](#gdenv-unpin)**
- **[gdenv vendor](#gdenv-vendor)**
- **[gdenv which](#gdenv-which)**

## **gdenv `completions`**

> ⚠️ **WARNING:** This command is not yet implemented.

Provides shell completions for the `gdenv` CLI application.

### Usage

`gdenv completions [OPTIONS] [SHELL]`

### Options

- `-o`, `--output <OUT_FILE>` — File to write the completions to
  - Default value: `stdout`

### Arguments

- `<SHELL>` — the specific version string to install (must be exact)
  - Supported values:
    - `bash`
    - `fish`
    - `powershell`
    - `zsh`

## **gdenv `install`**

Downloads and caches a specific version of _Godot_. If `VERSION` is omitted then the version is resolved using `-g`, `-p`, or `$PWD`.

### Usage

`gdenv install [OPTIONS] [VERSION]`

### Options

- `-f`, `--force` — forcibly overwrite an existing cache entry
- `-g`, `--global` — update the global pin (if `VERSION` is specified) or resolve `VERSION` from the global pin
- `-p`, `--path <PATH>` — resolve the pinned `VERSION` at `PATH`
- `-s`, `--src`, `--source` — install source code instead of an executable (cannot be used with `-g`)

### Arguments

- `[VERSION]` — the specific version string to install (must be exact)
  - Default value: Resolves the pinned version using `-g`, `-p`, or `$PWD` (if `-p` and `-g` omitted)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

## **gdenv `ls`/`list`**

Prints the path and version of all of the installed versions of _Godot_.

### Usage

`gdenv ls [OPTIONS]`

### Options

- `-a`, `--all` — list executable _and_ source code versions
- `-s`, `--src`, `--source` — list source code versions

## **gdenv `pin`**

Sets the _Godot_ version globally or for a specific directory.

### Usage

`gdenv pin [OPTIONS] <VERSION>`

### Options

- `-g`, `--global` — pin the system version (cannot be used with `-p`)
- `-i`, `--install` — installs the specified version of _Godot_ if missing
- `-f`, `--force` — forcibly overwrite an existing cache entry (only used with `-i`)
- `-p`, `--path <PATH>` — pin the specified path (cannot be used with `-g`)
  - Default value: `$PWD` (current working directory)

### Arguments

- `<VERSION>` — the specific version string to install (must be exact)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

## **gdenv `uninstall`**

Removes the specified version of _Godot_ from the `gdenv` download cache.

### Usage

`gdenv uninstall [OPTIONS] [VERSION]`

### Options

- `-a`, `--all` — uninstall all versions of _Godot_ (ignores source code without `-s`)
- `-s`, `--src`, `--source` — uninstall source code versions

### Arguments

- `[VERSION]` — the specific version string to install (must be exact; omit if using `-a`)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

## **gdenv `unpin`**

Removes a `Godot` version pin from the system or specified directory.

### Usage

`gdenv unpin [OPTIONS]`

### Options

- `-g`, `--global` — unpin the system version (cannot be used with `-p`)
- `-p`, `--path <PATH>` — unpin the specified path (cannot be used with `-g`)
  - Default value: `$PWD` (current working directory)

## **gdenv `vendor`**

Download the _Godot_ source code to the specified directory.

### Usage

`gdenv vendor [OPTIONS] [VERSION]`

### Options

- `-f`, `--force` — forcibly overwrite an existing cache entry
- `-o`, `--out <OUT_DIR>` — directory to extract the source code into (overwrites conflicting files)
  - Default value: `$PWD/godot-<VERSION>`
- `-p`, `--path <PATH>` — resolve the pinned `VERSION` at `PATH`
  - Default value: `$PWD` (current working directory)

### Arguments

- `[VERSION]` — the specific version string to install (must be exact and cannot be used with `-p`)
  - Default value: Resolves the pinned version at `$PWD`
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`

## **gdenv `which`**

Prints the path to the _Godot_ executable which would be used in the specified directory.

### Usage

`gdenv which [OPTIONS]`

### Options

- `-p`, `--path <PATH>` — the specified path to check
  - Default value: `$PWD` (current working directory)
