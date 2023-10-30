# Commands

## **gdenv `pin`**

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

## **gdenv `unpin`**

Removes a `Godot` version pin from the system or specified directory.

**Options:**

- **`-g`**, **`--global`** — unpin the system version (cannot be used with `-p`)
- **`-p`**, **`--path <PATH>`** — unpin the specified path (cannot be used with `-g`)
  - Default value: `$PWD` (current working directory)

## **gdenv `install`**

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

## **gdenv `vendor`**

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

## **gdenv `uninstall`**

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

## **gdenv `ls`/`list`**

Prints the path and version of all of the installed versions of _Godot_.

**Options:**

- **`-a`**, **`--all`** — list executable _and_ source code versions
- **`-s`**, **`--src`**, **`--source`** — list source code versions

## **gdenv `which`**

Prints the path to the _Godot_ executable which would be used in the specified directory.

**Options:**

- **`-p`**, **`--path <PATH>`** — the specified path to check
  - Default value: `$PWD` (current working directory)

## **gdenv `completions`**

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
