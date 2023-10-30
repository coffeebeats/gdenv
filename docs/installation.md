
# **Installation**

The easiest way to install `gdenv` is by using the pre-built binaries. These can be manually downloaded and configured, but automated installation scripts are provided and recommended.

> ⚠️ **WARNING:** It's good practice to inspect an installation script prior to execution. The scripts are included in this repository and can be reviewed prior to use.

## **Linux/MacOS — `sh` (recommended)**

```sh
curl https://raw.githubusercontent.com/coffeebeats/gdenv/main/scripts/install.sh | sh
```

## **Windows - `Git BASH for Windows` (recommended)**

> ❕ **NOTE:** If you're using [Git BASH for Windows](https://gitforwindows.org/), use the [Linux/MacOS instructions](#linuxmacos--sh-recommended).

## **Windows — `powershell` (recommended)**

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

## **Manual (not recommended)**

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

## **Install from source - Go (not recommended)**

`gdenv` is a Go project and can be installed using `go install`. This option is not recommended as it requires having the Go toolchain installed, it's slower than downloading a prebuilt binary, and there may be instability due to using a different version of Go than it was developed with.

> NOTE: You will need to somehow set the installed `gdenv-shim` binary as your system's `godot` command (consider using a symbolic link). This is done automatically by the recommended installation methods listed above.

```sh
go install github.com/coffeebeats/gdenv/cmd/gdenv@latest
```
