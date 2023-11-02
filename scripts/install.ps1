#requires -version 2

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
# https://github.com/coffeebeats/gdenv/blob/v0.5.2/docs/installation.md # x-release-please-version

<#
.SYNOPSIS
  Install 'gdenv' for managing multiple versions of the Godot editor.

.DESCRIPTION
  Install 'gdenv' for managing multiple versions of the Godot editor.

  The following environment variables can be set to modify behavior:

    - GDENV_VERSION: Install the specified version of 'gdenv'.

.PARAMETER NoModifyPath
    Do not modify the \$PATH environment variable.

.INPUTS
  <Inputs if any, otherwise state None>

.OUTPUTS
  <Outputs if any, otherwise state None - example: Log file stored in C:\Windows\Temp\<name>.log>

.NOTES
  Version:        0.5.2 # x-release-please-version
  Author:         https://github.com/coffeebeats
#>

# ------------------------------ Define: Params ------------------------------ #

param (
  [Switch] $NoModifyPath = $False
)
  
# ------------------------------ Define: Version ----------------------------- #
  
$GdEnvVersion = Get-GdEnvVersion

$GdEnvArchive = "gdenv-${GdEnvVersion}-windows-x86_64.zip"
  
# ------------------------------- Define: Store ------------------------------ #

$GdEnvHome = Get-GdEnvHome

Write-Host "using `$GDENV_HOME: ${GdEnvHome}"
  
# ----------------------------- Function: Install ---------------------------- #
  
$GdEnvRepositoryUrl = "https://github.com/coffeebeats/gdenv"

Function Install() {
  $GdEnvTempFolder = New-TemporaryFolder

  $GdEnvArchiveURL = "${$GdEnvRepositoryUrl}/releases/download/${GdEnvVersion}/${GdEnvArchive}"
  $GdEnvDownloadTo = "${GdEnvTempFolder}\${GdEnvArchive}"

  try {
    Invoke-WebRequest -URI $GdEnvArchiveURL -OutFile $GdEnvDownloadTo

    Microsoft.PowerShell.Archive\Expand-Archive `
    -Force `
    -Path $GdEnvDownloadTo `
    -DestinationPath "${GdEnvHome}\bin"

    [System.Environment]::SetEnvironmentVariable("GDENV_HOME", $GdEnvHome, "User")
  
    if (!($NoModifyPath)) {
      $GdEnvBinPath = "${$GdEnvHome}\bin"
    
      $PathParts = [System.Environment]::GetEnvironmentVariable("PATH", "User") -Split ";"
      $PathParts = $PathParts | Where-Object -ne $GdEnvBinPath
      $PathParts = $PathParts + $GdEnvBinPath

      [System.Environment]::SetEnvironmentVariable("PATH", $($PathParts -Join ";"), "User")
    }
  } catch {
    Write-Host "failed to install 'gdenv'"
  } finally {
    Remove-Item -Recurse $GdEnvTempFolder
  }
}

Install

# -------------------------- Function: Get-GdEnvHome ------------------------- #

# Returns the current value of the 'GDENV_HOME' environment variable or a
# default if unset.
Function Get-GdEnvHome() {
  if ([string]::IsNullOrEmpty($env:GDENV_HOME)) {
    return "${env:LOCALAPPDATA}\gdenv"
  }

  return $env:GDENV_HOME
}

# ------------------------ Function: Get-GdEnvVersion ------------------------ #

Function Get-GdEnvVersion() {
  if ([string]::IsNullOrEmpty($env:GDENV_VERSION)) {
    return "v0.5.2" # x-release-please-version
  }

  $GdEnvVersion = $env:GDENV_VERSION
  return "v" + $GdEnvVersion.TrimStart("v")
}

# --------------------- Function: Create-Temporary-Folder -------------------- #

# Creates a new temporary directory for downloading and extracting 'gdenv'.
Function New-TemporaryFolder() {
  # Make a new temporary folder with a randomized suffix
  $Name=$([System.IO.Path]::GetFileNameWithoutExtension([System.IO.Path]::GetRandomFileName()))
  $TemporaryFolderPath="${env:temp}\gdenv-${Name}"
  
  New-Item -ItemType Directory -Path $TemporaryFolderPath

  return $TemporaryFolderPath
}