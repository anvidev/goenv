package goenv

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var keys []string = []string{"TEST_STRING", "TEST_INT", "TEST_BOOL", "TEST_EXISTING", "TEST_MISSING",
	"TEST_TAGGED", "DB_HOST", "DB_PORT", "APP_NAME", "TEST_INT8", "TEST_INT64",
	"TEST_UINT", "TEST_FLOAT32", "TEST_FLOAT64", "TEST_INVALID_INT",
	"TEST_BOOL_1", "TEST_BOOL_0", "TEST_BOOL_TRUE", "TEST_BOOL_FALSE",
	"TEST_EMPTY", "LEVEL1_VAR", "LEVEL2_VAR", "LEVEL3_VAR", "TEST_DURATION", "TEST_TIME", "TEST_MISSING_DEFAULT",
	"SERVER_VERSION", "SERVER_ENV", "SERVER_ADDR", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT",
	"AUTH_TOKEN", "DATABASE_URL", "MAILER_API_KEY",
}

func TestStruct(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		input any
		fail  bool
		check func(t *testing.T, input any)
	}{
		{
			name: "Simple struct with all fields set",
			setup: func() {
				os.Setenv("TEST_STRING", "hello")
				os.Setenv("TEST_INT", "42")
				os.Setenv("TEST_BOOL", "true")
			},
			input: &struct {
				StringField string `goenv:"TEST_STRING"`
				IntField    int    `goenv:"TEST_INT"`
				BoolField   bool   `goenv:"TEST_BOOL"`
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					StringField string `goenv:"TEST_STRING"`
					IntField    int    `goenv:"TEST_INT"`
					BoolField   bool   `goenv:"TEST_BOOL"`
				})
				if s.StringField != "hello" {
					t.Errorf("StringField = %v, want %v", s.StringField, "hello")
				}
				if s.IntField != 42 {
					t.Errorf("IntField = %v, want %v", s.IntField, 42)
				}
				if s.BoolField != true {
					t.Errorf("BoolField = %v, want %v", s.BoolField, true)
				}
			},
		},
		{
			name: "Struct with missing optional environment variable",
			setup: func() {
				os.Setenv("TEST_EXISTING", "value")
			},
			input: &struct {
				ExistingField string `goenv:"TEST_EXISTING"`
				MissingField  string `goenv:"TEST_MISSING"`
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					ExistingField string `goenv:"TEST_EXISTING"`
					MissingField  string `goenv:"TEST_MISSING"`
				})
				if s.ExistingField != "value" {
					t.Errorf("ExistingField = %v, want %v", s.ExistingField, "value")
				}
				if s.MissingField != "" {
					t.Errorf("MissingFiled = %v, want %v", s.MissingField, "")
				}
			},
		},
		{
			name: "Struct with missing optional environment variable with a default",
			setup: func() {
				os.Setenv("TEST_EXISTING", "value")
			},
			input: &struct {
				ExistingField string `goenv:"TEST_EXISTING"`
				MissingField  string `goenv:"TEST_MISSING_DEFAULT,default=defaultvalue"`
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					ExistingField string `goenv:"TEST_EXISTING"`
					MissingField  string `goenv:"TEST_MISSING_DEFAULT,default=defaultvalue"`
				})
				fmt.Printf("input is: %+v\n", s)
				if s.ExistingField != "value" {
					t.Errorf("ExistingField = %v, want %v", s.ExistingField, "value")
				}
				if s.MissingField != "defaultvalue" {
					t.Errorf("MissingField = %v, want %v", s.MissingField, "defaultvalue")
				}
			},
		},
		{
			name: "Struct with missing required environment variable",
			setup: func() {
				os.Setenv("TEST_EXISTING", "value")
			},
			input: &struct {
				ExistingField string `goenv:"TEST_EXISTING"`
				MissingField  string `goenv:"TEST_MISSING,required"`
			}{},
			fail: true,
		},
		{
			name: "Struct with fields without goenv tags (should be skipped)",
			setup: func() {
				os.Setenv("TEST_TAGGED", "tagged_value")
			},
			input: &struct {
				TaggedField   string `goenv:"TEST_TAGGED"`
				UntaggedField string
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					TaggedField   string `goenv:"TEST_TAGGED"`
					UntaggedField string
				})
				if s.TaggedField != "tagged_value" {
					t.Errorf("TaggedField = %v, want %v", s.TaggedField, "tagged_value")
				}
				if s.UntaggedField != "" {
					t.Errorf("UntaggedField = %v, want %v", s.UntaggedField, "")
				}
			},
		},
		{
			name: "Nested struct",
			setup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_PORT", "5432")
				os.Setenv("APP_NAME", "myapp")
			},
			input: &struct {
				Database struct {
					Host string `goenv:"DB_HOST"`
					Port int    `goenv:"DB_PORT"`
				}
				AppName string `goenv:"APP_NAME"`
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					Database struct {
						Host string `goenv:"DB_HOST"`
						Port int    `goenv:"DB_PORT"`
					}
					AppName string `goenv:"APP_NAME"`
				})
				if s.Database.Host != "localhost" {
					t.Errorf("Database.Host = %v, want %v", s.Database.Host, "localhost")
				}
				if s.Database.Port != 5432 {
					t.Errorf("Database.Port = %v, want %v", s.Database.Port, 5432)
				}
				if s.AppName != "myapp" {
					t.Errorf("AppName = %v, want %v", s.AppName, "myapp")
				}
			},
		},
		{
			name:  "Invalid input - not a pointer",
			setup: func() {},
			input: struct {
				Field string `goenv:"TEST"`
			}{},
			fail: true,
		},
		{
			name:  "Invalid input - pointer to non-struct",
			setup: func() {},
			input: func() *string {
				s := "not a struct"
				return &s
			}(),
			fail: true,
		},
		{
			name: "Different data types",
			setup: func() {
				os.Setenv("TEST_INT8", "127")
				os.Setenv("TEST_INT64", "9223372036854775807")
				os.Setenv("TEST_UINT", "42")
				os.Setenv("TEST_FLOAT32", "3.14")
				os.Setenv("TEST_FLOAT64", "2.718281828")
				os.Setenv("TEST_DURATION", "5h")
				os.Setenv("TEST_TIME", "1992-06-26")
			},
			input: &struct {
				Int8Field     int8          `goenv:"TEST_INT8"`
				Int64Field    int64         `goenv:"TEST_INT64"`
				UintField     uint          `goenv:"TEST_UINT"`
				Float32Field  float32       `goenv:"TEST_FLOAT32"`
				Float64Field  float64       `goenv:"TEST_FLOAT64"`
				DurationField time.Duration `goenv:"TEST_DURATION"`
				TimeField     time.Time     `goenv:"TEST_TIME"`
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					Int8Field     int8          `goenv:"TEST_INT8"`
					Int64Field    int64         `goenv:"TEST_INT64"`
					UintField     uint          `goenv:"TEST_UINT"`
					Float32Field  float32       `goenv:"TEST_FLOAT32"`
					Float64Field  float64       `goenv:"TEST_FLOAT64"`
					DurationField time.Duration `goenv:"TEST_DURATION"`
					TimeField     time.Time     `goenv:"TEST_TIME"`
				})
				if s.Int8Field != 127 {
					t.Errorf("Int8Field = %v, want %v", s.Int8Field, 127)
				}
				if s.Int64Field != 9223372036854775807 {
					t.Errorf("Int64Field = %v, want %v", s.Int64Field, int64(9223372036854775807))
				}
				if s.UintField != 42 {
					t.Errorf("UintField = %v, want %v", s.UintField, uint(42))
				}
				if s.Float32Field != 3.14 {
					t.Errorf("Float32Field = %v, want %v", s.Float32Field, float32(3.14))
				}
				if s.Float64Field != 2.718281828 {
					t.Errorf("Float64Field = %v, want %v", s.Float64Field, 2.718281828)
				}
				if s.TimeField != time.Date(1992, 06, 26, 0, 0, 0, 0, time.UTC) {
					t.Errorf("TimeField = %v, want %v", s.TimeField, time.Date(1992, 06, 26, 0, 0, 0, 0, time.UTC))
				}
			},
		},
		{
			name: "Invalid type conversion",
			setup: func() {
				os.Setenv("TEST_INVALID_INT", "not_a_number")
			},
			input: &struct {
				IntField int `goenv:"TEST_INVALID_INT"`
			}{},
			fail: true,
		},
		{
			name: "Boolean edge cases",
			setup: func() {
				os.Setenv("TEST_BOOL_1", "1")
				os.Setenv("TEST_BOOL_0", "0")
				os.Setenv("TEST_BOOL_TRUE", "true")
				os.Setenv("TEST_BOOL_FALSE", "false")
			},
			input: &struct {
				Bool1     bool `goenv:"TEST_BOOL_1"`
				Bool0     bool `goenv:"TEST_BOOL_0"`
				BoolTrue  bool `goenv:"TEST_BOOL_TRUE"`
				BoolFalse bool `goenv:"TEST_BOOL_FALSE"`
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					Bool1     bool `goenv:"TEST_BOOL_1"`
					Bool0     bool `goenv:"TEST_BOOL_0"`
					BoolTrue  bool `goenv:"TEST_BOOL_TRUE"`
					BoolFalse bool `goenv:"TEST_BOOL_FALSE"`
				})
				if s.Bool1 != true {
					t.Errorf("Bool1 = %v, want %v", s.Bool1, true)
				}
				if s.Bool0 != false {
					t.Errorf("Bool0 = %v, want %v", s.Bool0, false)
				}
				if s.BoolTrue != true {
					t.Errorf("BoolTrue = %v, want %v", s.BoolTrue, true)
				}
				if s.BoolFalse != false {
					t.Errorf("BoolFalse = %v, want %v", s.BoolFalse, false)
				}
			},
		},
		{
			name: "Empty string environment variable",
			setup: func() {
				os.Setenv("TEST_EMPTY", "")
			},
			input: &struct {
				EmptyField string `goenv:"TEST_EMPTY"`
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					EmptyField string `goenv:"TEST_EMPTY"`
				})
				if s.EmptyField != "" {
					t.Errorf("EmptyField = %v, want %v", s.EmptyField, "")
				}
			},
		},
		{
			name: "Deeply nested struct",
			setup: func() {
				os.Setenv("LEVEL1_VAR", "level1")
				os.Setenv("LEVEL2_VAR", "level2")
				os.Setenv("LEVEL3_VAR", "level3")
			},
			input: &struct {
				Level1 struct {
					Var    string `goenv:"LEVEL1_VAR"`
					Level2 struct {
						Var    string `goenv:"LEVEL2_VAR"`
						Level3 struct {
							Var string `goenv:"LEVEL3_VAR"`
						}
					}
				}
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					Level1 struct {
						Var    string `goenv:"LEVEL1_VAR"`
						Level2 struct {
							Var    string `goenv:"LEVEL2_VAR"`
							Level3 struct {
								Var string `goenv:"LEVEL3_VAR"`
							}
						}
					}
				})
				if s.Level1.Var != "level1" {
					t.Errorf("Level1.Var = %v, want %v", s.Level1.Var, "level1")
				}
				if s.Level1.Level2.Var != "level2" {
					t.Errorf("Level1.Level2.Var = %v, want %v", s.Level1.Level2.Var, "level2")
				}
				if s.Level1.Level2.Level3.Var != "level3" {
					t.Errorf("Level1.Level2.Level3.Var = %v, want %v", s.Level1.Level2.Level3.Var, "level3")
				}
			},
		},
		{
			name: "Complex config struct",
			setup: func() {
				os.Setenv("SERVER_VERSION", "")
				os.Setenv("SERVER_ENV", "")
				os.Setenv("SERVER_ADDR", "")
				os.Setenv("SERVER_READ_TIMEOUT", "")
				os.Setenv("SERVER_WRITE_TIMEOUT", "")
				os.Setenv("SERVER_IDLE_TIMEOUT", "")
				os.Setenv("AUTH_TOKEN", "abc123")
				os.Setenv("DATABASE_URL", "file://database.db")
				os.Setenv("MAILER_API_KEY", "foobarbaz")
			},
			input: &struct {
				Server struct {
					Version      string        `goenv:"SERVER_VERSION,default=v0.1.0"`
					Env          string        `goenv:"SERVER_ENV,default=development"`
					Addr         string        `goenv:"SERVER_ADDR,default=:9090"`
					ReadTimeout  time.Duration `goenv:"SERVER_READ_TIMEOUT,default=10s"`
					WriteTimeout time.Duration `goenv:"SERVER_WRITE_TIMEOUT,default=30s"`
					IdleTimeout  time.Duration `goenv:"SERVER_IDLE_TIMEOUT,default=1m"`
				}
				Database struct {
					Token string `goenv:"AUTH_TOKEN,required"`
					URL   string `goenv:"DATABASE_URL,required"`
				}
				Mailer struct {
					From   string `goenv:"MAILER_FROM,default=Appname <noreply@appname.com>"`
					ApiKey string `goenv:"MAILER_API_KEY,required"`
				}
			}{},
			fail: false,
			check: func(t *testing.T, input any) {
				s := input.(*struct {
					Server struct {
						Version      string        `goenv:"SERVER_VERSION,default=v0.1.0"`
						Env          string        `goenv:"SERVER_ENV,default=development"`
						Addr         string        `goenv:"SERVER_ADDR,default=:9090"`
						ReadTimeout  time.Duration `goenv:"SERVER_READ_TIMEOUT,default=10s"`
						WriteTimeout time.Duration `goenv:"SERVER_WRITE_TIMEOUT,default=30s"`
						IdleTimeout  time.Duration `goenv:"SERVER_IDLE_TIMEOUT,default=1m"`
					}
					Database struct {
						Token string `goenv:"AUTH_TOKEN,required"`
						URL   string `goenv:"DATABASE_URL,required"`
					}
					Mailer struct {
						From   string `goenv:"MAILER_FROM,default=Appname <noreply@appname.com>"`
						ApiKey string `goenv:"MAILER_API_KEY,required"`
					}
				})

				if s.Server.Version != "v0.1.0" {
					t.Errorf("Server.Version = %v, want %v", s.Server.Version, "v0.1.0")
				}
				if s.Server.Env != "development" {
					t.Errorf("Server.Env = %v, want %v", s.Server.Env, "development")
				}
				if s.Server.Addr != ":9090" {
					t.Errorf("Server.Addr = %v, want %v", s.Server.Addr, ":9090")
				}
				if s.Server.ReadTimeout != 10*time.Second {
					t.Errorf("Server.ReadTimeout = %v, want %v", s.Server.ReadTimeout, 10*time.Second)
				}
				if s.Server.WriteTimeout != 30*time.Second {
					t.Errorf("Server.WriteTimeout = %v, want %v", s.Server.WriteTimeout, 30*time.Second)
				}
				if s.Server.IdleTimeout != 60*time.Second {
					t.Errorf("Server.IdleTimeout = %v, want %v", s.Server.IdleTimeout, 60*time.Second)
				}
				if s.Database.Token != "abc123" {
					t.Errorf("Database.Token = %v, want %v", s.Database.Token, "abc123")
				}
				if s.Database.URL != "file://database.db" {
					t.Errorf("Database.Token = %v, want %v", s.Database.URL, "file://database.db")
				}
				if s.Mailer.From != "Appname <noreply@appname.com>" {
					t.Errorf("Database.From = %v, want %v", s.Mailer.From, "Appname <noreply@appname.com>")
				}
				if s.Mailer.ApiKey != "foobarbaz" {
					t.Errorf("Database.ApiKey = %v, want %v", s.Mailer.ApiKey, "foobarbaz")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, key := range keys {
				// unset keys before tests to be sure
				os.Unsetenv(key)
			}

			if tt.setup != nil {
				tt.setup()
			}

			t.Cleanup(func() {
				for _, key := range keys {
					os.Unsetenv(key)
				}
			})

			err := Struct(tt.input)

			if tt.fail {
				if err == nil {
					t.Errorf("Struct() error = nil, fail %v", tt.fail)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Struct() error = %v, fail %v", err, tt.fail)
				return
			}

			if tt.check != nil {
				tt.check(t, tt.input)
			}
		})
	}
}
