# goenv

Simple Go package for retrieving environment variables with fallback values and type conversion.

## Installation

```bash
go get github.com/anvidev/goenv
```

## Basic Usage

> [!NOTE]
> Make sure to load your environment variables. See section [Loading environment variables](#loading-environment-variables) for more.

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

## Loading environment variables

To load your environment variables, simply place the following code in your main function.

```go

package main

import (
    // other imports

    "github.com/anvidev/goenv"

    // more imports
)

func main() {
    err := goenv.Load()
    if err != nil {
        // handle error
    }

    // your code goes here
}

```

> [!CAUTION]
> Running `Load` multiple times with different files might override values for duplicate keys

`Load` will by default try to load ".env".

However you can also provide file path(s) to `Load`, to load multiple files or from a specific location.

The code would then look like this:

```go
    err := goenv.Load("path/to/environment")
    if err != nil {
        // handle error
    }
```

## Functions

- `String(key, fallback string) string` - Get string with fallback
- `Int(key string, fallback int) int` - Get integer with fallback
- `Bool(key string, fallback bool) bool` - Get boolean with fallback
- `Duration(key string, fallback time.Duration) time.Duration` - Get duration with fallback
- `MustString(key string) string` - Get required string (panics if empty/unset)
- `Load(filenames ...string) error` - Loads 1 or more files in the environment. If no file is provided ".env" is used.

All functions (except `Load` and `MustString`) return the fallback value if the environment variable is not set or cannot be parsed.

## License

MIT
