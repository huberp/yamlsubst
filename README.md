# yamlsubst

[![CI](https://github.com/huberp/yamlsubst/actions/workflows/ci.yml/badge.svg)](https://github.com/huberp/yamlsubst/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/huberp/yamlsubst)](https://goreportcard.com/report/github.com/huberp/yamlsubst)
[![codecov](https://codecov.io/gh/huberp/yamlsubst/branch/main/graph/badge.svg)](https://codecov.io/gh/huberp/yamlsubst)
[![License](https://img.shields.io/github/license/huberp/yamlsubst)](LICENSE)
[![Release](https://img.shields.io/github/v/release/huberp/yamlsubst)](https://github.com/huberp/yamlsubst/releases)

A command-line tool for replacing placeholders in text with values from YAML files, similar to `envsubst` but powered by YAML.

## Features

- Replace placeholders in text files or stdin with values from YAML files
- Support for nested YAML paths (e.g., `${.app.config.host}`)
- Cross-platform support (Windows, Linux)
- Cross-architecture support (AMD64, ARM64)
- Simple and intuitive CLI interface

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/huberp/yamlsubst/releases).

### Build from Source

```bash
git clone https://github.com/huberp/yamlsubst.git
cd yamlsubst
./scripts/build.sh
```

Or on Windows:
```powershell
./scripts/build.ps1
```

### Install with Go

```bash
go install github.com/huberp/yamlsubst/cmd/yamlsubst@latest
```

## Usage

### Basic Usage

Create a YAML file with your values:

```yaml
# values.yaml
name: John Doe
age: 30
location:
  city: Seattle
  state: WA
```

Use placeholders in your template:

```
Hello, my name is ${.name} and I am ${.age} years old.
I live in ${.location.city}, ${.location.state}.
```

Process the template:

```bash
# From stdin
echo "Hello, my name is \${.name}" | yamlsubst --yaml values.yaml

# From file
yamlsubst --yaml values.yaml --file template.txt
```

Output:
```
Hello, my name is John Doe and I am 30 years old.
I live in Seattle, WA.
```

### Placeholder Syntax

Placeholders use the format `${.path.to.value}` where the path is a dot-separated sequence of keys to navigate the YAML structure.

Examples:
- `${.name}` - top-level key
- `${.person.name}` - nested key
- `${.app.config.host}` - deeply nested key

### Arithmetic Expressions

yamlsubst supports arithmetic expressions inside placeholders, allowing you to perform calculations using values from YAML files, environment variables, and literal numbers.

#### Supported Operations

- **Addition**: `+`
- **Subtraction**: `-`
- **Multiplication**: `*`
- **Division**: `/`
- **Parentheses**: `()` for grouping and controlling precedence

#### Operator Precedence

Expressions follow standard mathematical precedence rules:
1. Parentheses (highest)
2. Multiplication and Division
3. Addition and Subtraction (lowest)

#### Expression Components

Expressions can contain:
- **Literal numbers**: integers and floats (e.g., `42`, `3.14`, `0.5`)
- **YAML references**: starting with a dot (e.g., `.width`, `.app.config.port`)
- **Environment variables**: starting with a dollar sign (e.g., `$PORT`, `$DATABASE_PORT`)

#### Examples

**Basic arithmetic with YAML values:**
```yaml
# config.yaml
width: 10
height: 20
```
```bash
echo "Area: \${.width * .height} square units" | yamlsubst --yaml config.yaml
# Output: Area: 200 square units
```

**Combining literals and references:**
```yaml
# pricing.yaml
base_price: 100
```
```bash
echo "Total with tax: \${.base_price * 1.15}" | yamlsubst --yaml pricing.yaml
# Output: Total with tax: 115
```

**Using parentheses:**
```yaml
# calc.yaml
base: 50
offset: 10
multiplier: 2
```
```bash
echo "Result: \${(.base + .offset) * .multiplier}" | yamlsubst --yaml calc.yaml
# Output: Result: 120
```

**Environment variables in expressions:**
```bash
export PORT=8080
echo "Next port: \${$PORT + 1}" | yamlsubst --yaml values.yaml
# Output: Next port: 8081
```

**Mixed YAML and environment variables:**
```yaml
# config.yaml
base_port: 3000
```
```bash
export OFFSET=100
echo "Port: \${.base_port + $OFFSET}" | yamlsubst --yaml config.yaml
# Output: Port: 3100
```

**Complex expressions:**
```yaml
# dimensions.yaml
length: 15.5
width: 8.25
depth: 3
unit_price: 12.50
```
```bash
# Calculate volume and cost
echo "Volume: \${.length * .width * .depth} cubic units" | yamlsubst --yaml dimensions.yaml
# Output: Volume: 383.625 cubic units

echo "Cost: \$\${.length * .width * .depth * .unit_price}" | yamlsubst --yaml dimensions.yaml
# Output: Cost: $4795.3125
```

**Temperature conversion:**
```yaml
# weather.yaml
fahrenheit: 68
```
```bash
echo "Temperature: \${(.fahrenheit - 32) * 5 / 9}°C" | yamlsubst --yaml weather.yaml
# Output: Temperature: 20°C
```

#### Important Notes

- Division by zero returns an error and leaves the placeholder unchanged
- All numeric values are processed as floating-point numbers
- Results are formatted intelligently: whole numbers display without decimals (e.g., `10` not `10.0`)
- Invalid expressions leave the placeholder unchanged
- Non-numeric YAML values in arithmetic expressions will cause the placeholder to remain unchanged

### Command-Line Options

```
Usage:
  yamlsubst [flags]

Flags:
      --file string   Input file containing placeholders (reads from stdin if not specified)
  -h, --help          help for yamlsubst
      --yaml string   YAML file containing values for substitution (required)
```

### Examples

#### Example 1: Environment Configuration

```yaml
# config.yaml
database:
  host: localhost
  port: 5432
  name: myapp
```

```bash
echo "postgresql://\${.database.host}:\${.database.port}/\${.database.name}" | yamlsubst --yaml config.yaml
# Output: postgresql://localhost:5432/myapp
```

#### Example 2: Docker Compose Template

```yaml
# env.yaml
app:
  version: 1.0.0
  port: 8080
```

```yaml
# docker-compose.template.yml
version: '3'
services:
  app:
    image: myapp:${.app.version}
    ports:
      - "${.app.port}:8080"
```

```bash
yamlsubst --yaml env.yaml --file docker-compose.template.yml > docker-compose.yml
```

## Development

### Prerequisites

- Go 1.25.3 or higher
- Git

### Building

```bash
# Build for current platform
go build -o yamlsubst ./cmd/yamlsubst

# Build for all platforms
./scripts/build.sh
```

### Testing

```bash
# Run tests
./scripts/test.sh

# Run tests on Windows
./scripts/test.ps1

# Run specific test
go test -v ./pkg/substitutor/...
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`./scripts/test.sh`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by `envsubst` from GNU gettext
- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- YAML parsing by [go-yaml](https://github.com/go-yaml/yaml)