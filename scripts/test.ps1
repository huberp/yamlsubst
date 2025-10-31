# Test script for yamlsubst
# PowerShell version

$ErrorActionPreference = "Stop"

Write-Host "Running Go tests..." -ForegroundColor Green
go test -v -race -coverprofile=coverage.out ./...

if ($LASTEXITCODE -ne 0) {
    Write-Error "Tests failed!"
    exit 1
}

Write-Host ""
Write-Host "Coverage report:" -ForegroundColor Green
go tool cover -func=coverage.out

Write-Host ""
Write-Host "Running integration tests..." -ForegroundColor Green

# Create temporary test directory
$TEMP_DIR = New-Item -ItemType Directory -Path (Join-Path $env:TEMP ([System.IO.Path]::GetRandomFileName()))

try {
    # Test 1: Basic substitution
    @"
name: TestUser
version: 1.0.0
"@ | Out-File -FilePath "$TEMP_DIR/values.yaml" -Encoding UTF8

    @"
Name: `${.name}
Version: `${.version}
"@ | Out-File -FilePath "$TEMP_DIR/template.txt" -Encoding UTF8

    # Build the binary
    go build -o "$TEMP_DIR/yamlsubst.exe" ./cmd/yamlsubst

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Build failed!"
        exit 1
    }

    $RESULT = & "$TEMP_DIR/yamlsubst.exe" --yaml "$TEMP_DIR/values.yaml" --file "$TEMP_DIR/template.txt"
    $EXPECTED = "Name: TestUser`r`nVersion: 1.0.0"

    if ($RESULT -ne $EXPECTED) {
        Write-Error "Integration test failed!`nExpected: $EXPECTED`nGot: $RESULT"
        exit 1
    }

    # Test 2: Nested values
    @"
app:
  name: MyApp
  config:
    host: localhost
    port: 8080
"@ | Out-File -FilePath "$TEMP_DIR/nested.yaml" -Encoding UTF8

    $RESULT = "Server: `${.app.config.host}:`${.app.config.port}" | & "$TEMP_DIR/yamlsubst.exe" --yaml "$TEMP_DIR/nested.yaml"
    $EXPECTED = "Server: localhost:8080"

    if ($RESULT -ne $EXPECTED) {
        Write-Error "Nested values integration test failed!`nExpected: $EXPECTED`nGot: $RESULT"
        exit 1
    }

    # Test 3: Version command
    $VERSION_OUTPUT = & "$TEMP_DIR/yamlsubst.exe" version
    if ($VERSION_OUTPUT -notmatch "yamlsubst version") {
        Write-Error "Version command test failed!`nGot: $VERSION_OUTPUT"
        exit 1
    }

    # Test 4: Help command
    $HELP_OUTPUT = & "$TEMP_DIR/yamlsubst.exe" --help
    if ($HELP_OUTPUT -notmatch "yamlsubst is a CLI tool") {
        Write-Error "Help command test failed!"
        exit 1
    }

    Write-Host ""
    Write-Host "All integration tests passed! ✓" -ForegroundColor Green

} finally {
    # Cleanup
    Remove-Item -Path $TEMP_DIR -Recurse -Force
}
