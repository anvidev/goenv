// Package goenv provides functions for retrieving values from environment variables with fallback values.
package goenv

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// String retrives the value of environment variable `k`. If no value is found, then the fallback value is returned.
func String(k, f string) string {
	v, found := os.LookupEnv(k)
	if !found || v == "" {
		return f
	}
	return v
}

// Duration retrieves the value of environment variable `k`. If the variable is present, then the value is parsed
// into `time.Duration`. Should this fail, then the fallback value is returned. If the variable is not present, then
// the fallback value is returned.
func Duration(k string, f time.Duration) time.Duration {
	v, found := os.LookupEnv(k)
	if !found {
		return f
	}
	dur, err := time.ParseDuration(v)
	if err != nil {
		return f
	}
	return dur
}

// Int retrieves the value of environment variable `k`. If the variable is present, then the value is parsed
// into type `int`. If the variable is not present, then the fallback is returned.
func Int(k string, f int) int {
	v, found := os.LookupEnv(k)
	if !found {
		return f
	}
	int, err := strconv.Atoi(v)
	if err != nil {
		return f
	}
	return int
}

// Bool retrieves the value of environment variable `k`. If the variable is present, then the value is parsed
// into type `bool`. If the variable is not present, then the fallback is returned.
func Bool(k string, f bool) bool {
	v, found := os.LookupEnv(k)
	if !found {
		return f
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return f
	}
	return b
}

// MustString retrives the value of environment variable `k`. If no value is found, then the program panics.
func MustString(k string) string {
	v := String(k, "")
	if v == "" {
		panic(fmt.Errorf("environment variable %s is not defined", k))
	}
	return v
}
