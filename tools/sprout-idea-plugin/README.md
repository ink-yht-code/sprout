# Sprout IDEA Plugin

GoLand/IntelliJ IDEA plugin for `.sprout` files, providing syntax highlighting, basic syntax validation, and code completion.

## Features

- **File Type Recognition**: Automatically recognizes `.sprout` files
- **Syntax Highlighting**: Keywords, HTTP methods, strings, comments, operators
- **Syntax Validation**: Basic error detection for arrow return types
- **Code Completion**: 
  - Keywords: `type`, `server`, `prefix`, `public`, `private`, `service`, `rpc`
  - HTTP Methods: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`

## Requirements

- IntelliJ IDEA 2023.2+ or GoLand 2023.2+
- JDK 17+

## Building

```bash
# Windows
.\gradlew buildPlugin

# Linux/macOS
./gradlew buildPlugin
```

The plugin zip will be in `build/distributions/`.

## Running for Development

```bash
# Windows
.\gradlew runIde

# Linux/macOS
./gradlew runIde
```

This launches a sandboxed IDE with the plugin installed.

## Installation

1. Build the plugin or download from releases
2. In GoLand/IDEA: `File` → `Settings` → `Plugins` → `⚙️` → `Install Plugin from Disk...`
3. Select the zip file from `build/distributions/`
4. Restart IDE

## .sprout File Example

```sprout
// Type definitions
type HelloReq {
    Name string `json:"name"`
}

type HelloResp {
    Message string `json:"message"`
}

// Server configuration
server {
    prefix: "/api"
}

// Public API endpoints
public service HelloService {
    GET /hello(HelloReq) => HelloResp
    POST /greet(HelloReq) => HelloResp
}

// Private RPC methods
private service InternalService {
    rpc Ping(PingReq) => PingResp
}
```

## Color Settings

Customize colors in: `Settings` → `Editor` → `Color Scheme` → `Sprout`
