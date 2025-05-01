# This script installs 'gdenv' by downloading prebuilt binaries from the
# project's GitHub releases page. By default the latest version is installed,
# but a different release can be used instead by setting $GDENV_VERSION.
#
# The script will set up a 'gdenv' cache at '%LOCALAPPDATA%/gdenv'. This
# behavior can be customized by setting '$GDENV_HOME' prior to running the
# script. Existing Godot artifacts cached in a 'gdenv' store won't be lost, but
# this script will overwrite any 'gdenv' binary artifacts in '$GDENV_HOME/bin'.
#
# NOTE: Unlike the 'install.sh' counterpart, this script exclusively installs
# the 'gdenv' binary for 64-bit Windows. If an alternative 'gdenv' binary is
# required, follow the documentation for an alternative means of installation:
# https://github.com/coffeebeats/gdenv/blob/v0.6.29/docs/installation.md # x-release-please-version

<#
.SYNOPSIS
  Install 'gdenv' for managing multiple versions of the Godot editor.

.DESCRIPTION
  This script downloads the specified version of 'gdenv' from GitHub, extracts
  its artifacts to the 'gdenv' store ('$GDENV_HOME' or a default path), and then
  updates environment variables as needed.

.PARAMETER NoModifyPath
  Do not modify the $PATH environment variable.

.PARAMETER Version
  Install the specified version of 'gdenv'.

.INPUTS
  None

.OUTPUTS
  $env:GDENV_HOME\bin\gdenv.exe
  $env:GDENV_HOME\bin\godot.exe

.NOTES
  Version:        0.6.29 # x-release-please-version
  Author:         https://github.com/coffeebeats

.LINK
  https://github.com/coffeebeats/gdenv
#>

# ------------------------------ Define: Params ------------------------------ #

Param (
  # NoModifyPath - if set, the user's $PATH variable won't be updated
  [Switch] $NoModifyPath = $False,

  # Version - override the specific version of 'gdenv' to install
  [String] $Version = "v0.6.29" # x-release-please-version
)

# -------------------------- Function: Get-GdEnvHome ------------------------- #

# Returns the current value of the 'GDENV_HOME' environment variable or a
# default if unset.
Function Get-GdEnvHome() {
  if ([string]::IsNullOrEmpty($env:GDENV_HOME)) {
    return Join-Path -Path $env:LOCALAPPDATA -ChildPath "gdenv"
  }

  return $env:GDENV_HOME
}

# ------------------------ Function: Get-GdEnvVersion ------------------------ #

Function Get-GdEnvVersion() {
  return "v" + $Version.TrimStart("v")
}

# --------------------- Function: Create-Temporary-Folder -------------------- #

# Creates a new temporary directory for downloading and extracting 'gdenv'. The
# returned directory path will have a randomized suffix.
Function New-TemporaryFolder() {
  # Make a new temporary folder with a randomized suffix.
  return New-Item `
    -ItemType Directory `
    -Name "gdenv-$([System.IO.Path]::GetFileNameWithoutExtension([System.IO.Path]::GetRandomFileName()))"`
    -Path $env:temp
}

# ------------------------------- Define: Store ------------------------------ #

$GdEnvHome = Get-GdEnvHome

Write-Host "info: setting 'GDENV_HOME' environment variable: ${GdEnvHome}"

[System.Environment]::SetEnvironmentVariable("GDENV_HOME", $GdEnvHome, "User")

# ------------------------------ Define: Version ----------------------------- #
  
$GdEnvVersion = Get-GdEnvVersion

$GdEnvArchive = "gdenv-${GdEnvVersion}-windows-x86_64.zip"

# ----------------------------- Execute: Install ----------------------------- #
  
$GdEnvRepositoryURL = "https://github.com/coffeebeats/gdenv"

# Install downloads 'gdenv' and extracts its binaries into the store. It also
# updates environment variables as needed.
Function Install() {
  $GdEnvTempFolder = New-TemporaryFolder

  $GdEnvArchiveURL = "${GdEnvRepositoryURL}/releases/download/${GdEnvVersion}/${GdEnvArchive}"
  $GdEnvDownloadTo = Join-Path -Path $GdEnvTempFolder -ChildPath $GdEnvArchive

  $GdEnvHomeBinPath = Join-Path -Path $GdEnvHome -ChildPath "bin"

  try {
    Write-Host "info: installing version: '${GdEnvVersion}'"

    Invoke-WebRequest -URI $GdEnvArchiveURL -OutFile $GdEnvDownloadTo

    Microsoft.PowerShell.Archive\Expand-Archive `
      -Force `
      -Path $GdEnvDownloadTo `
      -DestinationPath $GdEnvHomeBinPath
  
    if (!($NoModifyPath)) {
      $PathParts = [System.Environment]::GetEnvironmentVariable("PATH", "User").Trim(";") -Split ";"
      $PathParts = $PathParts.where{ $_ -ne $GdEnvHomeBinPath }
      $PathParts = $PathParts + $GdEnvHomeBinPath

      Write-Host "info: updating 'PATH' environment variable: ${GdEnvHomeBinPath}"

      [System.Environment]::SetEnvironmentVariable("PATH", $($PathParts -Join ";"), "User")
    }

    Write-Host "info: sucessfully installed executables:`n"
    Write-Host "  gdenv.exe: $(Join-Path -Path $GdEnvHomeBinPath -ChildPath "gdenv.exe")"
    Write-Host "  godot.exe: $(Join-Path -Path $GdEnvHomeBinPath -ChildPath "godot.exe")"
  }
  catch {
    Write-Host "error: failed to install 'gdenv': ${_}"
  }
  finally {
    Write-Host "`ninfo: cleaning up downloads: ${GdEnvTempFolder}"

    Remove-Item -Recurse $GdEnvTempFolder
  }
}

Install
