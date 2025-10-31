# Build script for yamlsubst - cross-platform compilation
# PowerShell version

param(
    [string]$Version = "dev",
    [string]$Commit = "",
    [string]$Date = ""
)

$ErrorActionPreference = "Stop"

if ($Commit -eq "") {
    try {
        $Commit = git rev-parse --short HEAD
    } catch {
        $Commit = "none"
    }
}

if ($Date -eq "") {
    $Date = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
}

$LDFLAGS = "-X main.version=$Version -X main.commit=$Commit -X main.date=$Date"

# Create dist directory
New-Item -ItemType Directory -Force -Path dist | Out-Null

Write-Host "Building yamlsubst $Version..." -ForegroundColor Green

# Build for different platforms
$platforms = @(
    @{OS="linux"; ARCH="amd64"},
    @{OS="linux"; ARCH="arm64"},
    @{OS="windows"; ARCH="amd64"},
    @{OS="windows"; ARCH="arm64"}
)

foreach ($platform in $platforms) {
    $GOOS = $platform.OS
    $GOARCH = $platform.ARCH
    
    $outputName = "yamlsubst-$GOOS-$GOARCH"
    if ($GOOS -eq "windows") {
        $outputName = "$outputName.exe"
    }
    
    Write-Host "Building for $GOOS/$GOARCH..." -ForegroundColor Cyan
    
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    
    & go build -ldflags $LDFLAGS -o "dist/$outputName" ./cmd/yamlsubst
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Build failed for $GOOS/$GOARCH"
        exit 1
    }
}

Write-Host "Build complete! Binaries are in dist/" -ForegroundColor Green
Get-ChildItem dist/ | Format-Table Name, Length
