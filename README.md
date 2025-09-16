# Apix

A CLI application that simplifies API testing and interaction.

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Build Status](https://img.shields.io/badge/Build-Passing-success?style=for-the-badge)](https://github.com/yourusername/apix/actions)

## Overview

Unlike curl's verbose syntax and complex flag management, Apix provides an intuitive interface for making HTTP requests, managing authentication, and handling responses with built-in formatting and error handling.

## Features

- 🚀 **Simple HTTP Methods**: Support for GET, POST, PUT, and DELETE requests
- 🎯 **Interactive Mode**: Built-in interactive CLI using the `huh` package
- 📝 **Clean Syntax**: More intuitive than curl for everyday API testing
- 🔧 **Built-in Formatting**: Automatic response formatting and error handling
- 🔐 **Authentication Support**: Easy authentication management
- ⚡ **Fast & Lightweight**: Built with Go for performance

## Installation

### From Source

```bash
go install github.com/Esa824/apix@latest
```

### Build Locally

```bash
git clone https://github.com/Esa824/apix.git
cd apix
go build -o apix
```

## Usage

### Command Mode

Execute HTTP requests directly from the command line:

```bash
# GET request
apix get https://api.example.com/users

# POST request with data
apix post https://api.example.com/users --data '{"name":"John","email":"john@example.com"}'

# PUT request
apix put https://api.example.com/users/1 --data '{"name":"Jane Doe"}'

# DELETE request
apix delete https://api.example.com/users/1
```

### Interactive Mode

Launch the interactive CLI interface:

```bash
apix --cli
```

This opens an intuitive interface where you can:
- Select HTTP methods
- Enter URLs and parameters
- Add headers and authentication
- View formatted responses

## Examples

### Basic GET Request
```bash
apix get https://jsonplaceholder.typicode.com/posts/1
```

### POST with JSON Data
```bash
apix post https://httpbin.org/post \
  --header "Content-Type: application/json" \
  --data '{"key": "value", "number": 42}'
```

### Authentication Example
```bash
apix get https://api.github.com/user \
  --header "Authorization: Bearer your-token-here"
```

## Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `get` | Make GET request | `apix get https://api.example.com/data` |
| `post` | Make POST request | `apix post https://api.example.com/data --data '{}'` |
| `put` | Make PUT request | `apix put https://api.example.com/data/1 --data '{}'` |
| `delete` | Make DELETE request | `apix delete https://api.example.com/data/1` |
| `--cli` | Launch interactive mode | `apix --cli` |

## Configuration

Apix supports various configuration options:

- **Headers**: Add custom headers to your requests
- **Authentication**: Support for Bearer tokens, Basic auth, and custom auth
- **Output Formatting**: JSON pretty-printing and response highlighting
- **Request Timeout**: Configurable timeout settings

## Development

### Prerequisites

- Go 1.19 or later
- Git

### Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Huh](https://github.com/charmbracelet/huh) - Interactive forms

### Building

```bash
# Build for current platform
make
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Roadmap

- [ ] Response caching
- [ ] Request history
- [ ] Environment variable support
- [ ] Configuration file support
- [ ] Plugin system
- [ ] GraphQL support
- [ ] WebSocket support

## Why Apix?

### vs. curl
- **Simpler syntax**: No need to remember complex curl flags
- **Interactive mode**: Built-in forms for easy request building
- **Better output**: Automatic JSON formatting and syntax highlighting

### vs. Postman
- **Lightweight**: No GUI overhead, perfect for CI/CD
- **Version control friendly**: Text-based configuration
- **Terminal native**: Fits naturally into developer workflows

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for the excellent CLI framework
- [Charm](https://charm.sh/) for the beautiful terminal UI components
- The Go community for inspiration and support
