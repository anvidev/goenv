# goenv

Simple Go package for retrieving environment variables with fallback values and type conversion.

## Installation

```bash
go get github.com/anvidev/goenv
```

## Usage

```go
package main

import (
    "time"
    "github.com/anvidev/goenv"
)

func main() {
    // Get environment variables with fallbacks
    host := goenv.String("HOST", "localhost")
    port := goenv.Int("PORT", 8080)
    debug := goenv.Bool("DEBUG", false)
    timeout := goenv.Duration("TIMEOUT", 30*time.Second)
    
    // Required variable (panics if not set)
    apiKey := goenv.MustString("API_KEY")
}

// Common pattern for configuration structs
type Config struct {
    Host    string
    Port    int
    Debug   bool
    Timeout time.Duration
    APIKey  string
}

func LoadConfig() Config {
    return Config{
        Host:    goenv.String("HOST", "localhost"),
        Port:    goenv.Int("PORT", 8080),
        Debug:   goenv.Bool("DEBUG", false),
        Timeout: goenv.Duration("TIMEOUT", 30*time.Second),
        APIKey:  goenv.MustString("API_KEY"),
    }
}
```

## Functions

- `String(key, fallback string) string` - Get string with fallback
- `Int(key string, fallback int) int` - Get integer with fallback  
- `Bool(key string, fallback bool) bool` - Get boolean with fallback
- `Duration(key string, fallback time.Duration) time.Duration` - Get duration with fallback
- `MustString(key string) string` - Get required string (panics if empty/unset)

All functions return the fallback value if the environment variable is not set or cannot be parsed.

## License

MIT
