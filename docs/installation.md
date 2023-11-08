
# **Installation**

The easiest way to install `gdenv` is by using the pre-built binaries. These can be manually downloaded and configured, but automated installation scripts are provided and recommended.

Alternatively, you can install `gdenv` from source using the latest supported version of [Go](https://go.dev/). See [Install from source](#install-from-source) for more details.

## **Pre-built binaries**

> ⚠️ **WARNING:** It's good practice to inspect an installation script prior to execution. The scripts are included in this repository and can be reviewed prior to use.

### **Linux/MacOS (recommended)**

```sh
curl https://raw.githubusercontent.com/coffeebeats/gdenv/main/scripts/install.sh | sh
```

### **Windows (recommended)**

#### **Git BASH for Windows**

If you're using [Git BASH for Windows](https://gitforwindows.org/) follow the recommended [Linux/MacOS](#linuxmacos-recommended) instructions.

#### **Powershell**

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

### **Manual download**

> ❕ **NOTE:** The instructions below provide `bash`-specific commands for a _Linux_-based system. While these won't work in _PowerShell_, the process will be similar.

1. Download a prebuilt binary from the corresponding GitHub release.

    ```sh
    # Set '$VERSION', '$OS', and '$ARCH' to the desired values.
    VERSION=0.0.0 OS=linux ARCH=x86_64; \
    curl -LO https://github.com/coffeebeats/gdenv/releases/download/v$VERSION/gdenv-$VERSION-$OS-$ARCH.tar.gz
    ```

2. Extract the downloaded archive.

    ```sh
    # Set '$GDENV_HOME' to the desired location (defaults to '$HOME/.gdenv' on Linux/MacOS).
    GDENV_HOME=$HOME/.gdenv; \
    mkdir -p $GDENV_HOME/bin && \
    tar -C $GDENV_HOME/bin -xf gdenv-$VERSION-$OS-$ARCH.tar.gz
    ```

3. Export the `GDENV_HOME` environment variable and add `$GDENV_HOME/bin` to `$PATH`.

    ```sh
    # In '.bashrc' or something similar ('$GDENV_HOME' can be customized).
    export GDENV_HOME="$HOME/.gdenv"
    export PATH="$GDENV_HOME/bin:$PATH"
    ```

## **Install from source**

`gdenv` is a Go project and can be installed using `go install`. This option is not recommended as it requires having the Go toolchain installed, it's slower than downloading a prebuilt binary, and there may be instability due to using a different version of Go than it was developed with.

```sh
go install github.com/coffeebeats/gdenv/cmd/gdenv@latest
go install github.com/coffeebeats/gdenv/cmd/gdenv-shim@latest
```

Once `gdenv` and `gdenv-shim` are installed a few things need to be configured. Follow the instructions below based on your operating system.

### **Linux/MacOS**

1. Export the `GDENV_HOME` environment variable and add `$GDENV_HOME/bin` to the `PATH` environment variable.

    Add the following to your shell's profile script/RC file:

    ```sh
    export GDENV_HOME="$HOME/.gdenv"
    export PATH="$GDENV_HOME/bin:$PATH"
    ```

2. When installing from source the `gdenv-shim` binary is _not_ renamed to `godot`; that is only done as part of the build process for the published binaries. To resolve this and allow the `godot` command to route to the correct version of _Godot_, add a link from the installed `gdenv-shim` binary to `godot`:

    > ❕ **NOTE:** Make sure to restart your terminal after the previous step so that any changes to `$GDENV_HOME` have been updated.

    ```sh
    test ! -z $GDENV_HOME && \
        command -v gdenv-shim >/dev/null 2>&1 && \
        ln -s $(which gdenv-shim) $GDENV_HOME/bin/godot
    ```

### **Windows (Powershell)**

1. Export the `GDENV_HOME` environment variable using the following:

    ```sh
    $GdEnvHomePath = "${env:LOCALAPPDATA}\gdenv" # Replace with whichever path you'd like.
    [System.Environment]::SetEnvironmentVariable("GDENV_HOME", $GdEnvHomePath, "User")
    ```

2. Add `$GDENV_HOME/bin` to your `PATH` environment variable:

    > ❕ **NOTE:** Make sure to restart your terminal after the previous step so that any changes to `$GDENV_HOME` have been updated.

    ```sh
    $PathParts = [System.Environment]::GetEnvironmentVariable("PATH", "User").Trim(";") -Split ";"
    $PathParts = $PathParts.where{ $_ -ne "${env:GDENV_HOME}\bin" }
    $PathParts = $PathParts + "${env:GDENV_HOME}\bin"

    [System.Environment]::SetEnvironmentVariable("PATH", $($PathParts -Join ";"), "User")
    ```

3. When installing from source the `gdenv-shim` binary is _not_ renamed to `godot`; that is only done as part of the build process for the published binaries. To resolve this and allow the `godot` command to route to the correct version of _Godot_, add a link from the installed `gdenv-shim` binary to `godot`:

    > ❕ **NOTE:** Make sure to restart your terminal after the previous step so that any changes to `$GDENV_HOME` and `$PATH` have been updated.

    ```sh
    New-Item "${env:GDENV_HOME}\bin" -ItemType Directory -ea 0
    New-Item -Path "${env:GDENV_HOME}\bin\godot.exe" -ItemType SymbolicLink -Value $((Get-Command gdenv-shim).Path)
    ```

    If you encounter an error `New-Item: Administrator privilege required for this operation` then re-run the command in step 3 from a terminal with Administrator privileges.
