# goenv

Simple Go package for retrieving environment variables with fallback values and type conversion.

## Installation

```bash
go get github.com/anvidev/goenv
```
## Functions

- `String(key, fallback string) string` - Get string with fallback
- `Int(key string, fallback int) int` - Get integer with fallback  
- `Bool(key string, fallback bool) bool` - Get boolean with fallback
- `Duration(key string, fallback time.Duration) time.Duration` - Get duration with fallback
- `MustString(key string) string` - Get required string (panics if empty/unset)
- `Struct(v any) error` - Populate a struct using `goenv` struct tags
- `Load(filenames ...string) error` - Loads 1 or more files in the environment. If no file is provided ".env" is used.

## Basic Usage

> [!NOTE]
> Make sure to load your environment variables. See section [Loading environment variables](#loading-environment-variables) for more.

Using the primitives.

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

Using `goenv` struct tags

```go
package main

import (
    "fmt"
    "time"
    "github.com/anvidev/goenv"
)

type serverConfig struct {
    Env         string  `goenv:"ENV,default=development"`
    Port        int     `goenv:"PORT,required"`
    ServerName  string  `goenv:"SERVER_NAME"` 
}

func main() {
    var config serverConfig

    if err := goenv.Struct(&config); err != nil {
        log.Fatal(err)
    }

    // Config is populated
    fmt.Println(config.Env)
}
```

## Struct tags

The `Struct()` function iterates through struct fields and populates them based on
the `goenv` struct tag.

Supported field types:

 - string
 - int, int8, int16, int32, int64
 - uint, uint8, uint16, uint32, uint64
 - float32, float64
 - bool
 - time.Duration
 - time.Time (uses Golang's time formats)
 - nested structs (processed recursively)

| Fields   | Description                                                          |
|----------|----------------------------------------------------------------------|
| default  | Sets the field value to default if environment variable is not found |
| required | Returns an error if environment variable is not found                |

## Loading environment variables

To load your environment variables, simply place the following code in your main function.

```go

package main

import (
    // imports

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


## License

MIT
