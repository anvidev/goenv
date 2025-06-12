package goenv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

// Loads the content of 1 or more files in to the current environment.
//
// If no files are provided, Load defaults to ".env".
//
// If multiple files are provided the first file is loaded fully,
// while only keys not already in the environment are loaded for the remaining files.
func Load(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = append(filenames, ".env")
	}

	for i, filename := range filenames {
		if err := loadFile(filename, i == 0); err != nil {
			return fmt.Errorf("goenv: Failed to load file '%s': %s", filename, err.Error())
		}
	}

	return nil
}

func loadFile(filename string, overload bool) error {
	src, err := readFile(filename)
	if err != nil {
		return err
	}

	fileMap, err := parseInput(src)
	if err != nil {
		return err
	}

	curEnv := os.Environ()
	for key, value := range fileMap {
		exists := slices.ContainsFunc(curEnv, func(val string) bool {
			curKey := strings.Split(val, "=")[0]
			return curKey == key
		})

		if !exists || overload {
			os.Setenv(key, value)
		}
	}

	return nil
}

func readFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, f)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
