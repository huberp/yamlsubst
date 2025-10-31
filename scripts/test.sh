#!/bin/bash
# Test script for yamlsubst

set -e

echo "Running Go tests..."
go test -v -race -coverprofile=coverage.out -coverpkg=./pkg/... ./pkg/...

echo ""
echo "Coverage report:"
go tool cover -func=coverage.out

echo ""
echo "Running integration tests..."

# Create temporary test files
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Test 1: Basic substitution
cat > "$TEMP_DIR/values.yaml" << 'EOF'
name: TestUser
version: 1.0.0
EOF

cat > "$TEMP_DIR/template.txt" << 'EOF'
Name: ${.name}
Version: ${.version}
EOF

# Build the binary
go build -o "$TEMP_DIR/yamlsubst" ./cmd/yamlsubst

RESULT=$("$TEMP_DIR/yamlsubst" --yaml "$TEMP_DIR/values.yaml" --file "$TEMP_DIR/template.txt")
EXPECTED="Name: TestUser
Version: 1.0.0"

if [ "$RESULT" != "$EXPECTED" ]; then
    echo "Integration test failed!"
    echo "Expected:"
    echo "$EXPECTED"
    echo "Got:"
    echo "$RESULT"
    exit 1
fi

# Test 2: Nested values
cat > "$TEMP_DIR/nested.yaml" << 'EOF'
app:
  name: MyApp
  config:
    host: localhost
    port: 8080
EOF

RESULT=$(echo "Server: \${.app.config.host}:\${.app.config.port}" | "$TEMP_DIR/yamlsubst" --yaml "$TEMP_DIR/nested.yaml")
EXPECTED="Server: localhost:8080"

if [ "$RESULT" != "$EXPECTED" ]; then
    echo "Nested values integration test failed!"
    echo "Expected: $EXPECTED"
    echo "Got: $RESULT"
    exit 1
fi

# Test 3: Version command
VERSION_OUTPUT=$("$TEMP_DIR/yamlsubst" version)
if [[ ! "$VERSION_OUTPUT" =~ "yamlsubst version" ]]; then
    echo "Version command test failed!"
    echo "Got: $VERSION_OUTPUT"
    exit 1
fi

# Test 4: Help command
HELP_OUTPUT=$("$TEMP_DIR/yamlsubst" --help)
if [[ ! "$HELP_OUTPUT" =~ "yamlsubst is a CLI tool" ]]; then
    echo "Help command test failed!"
    exit 1
fi

echo ""
echo "All integration tests passed! ✓"
