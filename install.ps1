$ErrorActionPreference = "Stop"

$Repo = "ma-cohen/code-cat"
$Binary = "ccat"
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { "$env:USERPROFILE\.local\bin" }

function Get-Arch {
  switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { "amd64" }
    "ARM64" { "arm64" }
    default { throw "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
  }
}

function Get-LatestVersion {
  $release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
  return $release.tag_name
}

$Arch = Get-Arch
$Version = Get-LatestVersion
$Archive = "${Binary}_windows_${Arch}.zip"
$Url = "https://github.com/$Repo/releases/download/$Version/$Archive"

Write-Host "Installing $Binary $Version (windows/$Arch)..."

$Tmp = Join-Path $env:TEMP ([System.Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Path $Tmp | Out-Null

try {
  $ZipPath = Join-Path $Tmp $Archive
  Invoke-WebRequest -Uri $Url -OutFile $ZipPath -UseBasicParsing
  Expand-Archive -Path $ZipPath -DestinationPath $Tmp

  if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
  }

  Copy-Item (Join-Path $Tmp "$Binary.exe") (Join-Path $InstallDir "$Binary.exe") -Force

  # Add to PATH for current session if not already present
  if ($env:PATH -notlike "*$InstallDir*") {
    $env:PATH = "$InstallDir;$env:PATH"
    # Persist for future sessions (user scope)
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$InstallDir*") {
      [Environment]::SetEnvironmentVariable("PATH", "$InstallDir;$currentPath", "User")
      Write-Host "Added $InstallDir to your PATH (restart terminal to take effect)"
    }
  }

  Write-Host "Installed: $(& "$InstallDir\$Binary.exe" --version)"
} finally {
  Remove-Item -Recurse -Force $Tmp
}
