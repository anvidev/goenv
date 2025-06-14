package goenv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func loadFiles(filenames []string) error {
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
	curEnvMap := make(map[string]bool)
	for _, val := range curEnv {
		curKey := strings.Split(val, "=")[0]
		curEnvMap[curKey] = true
	}

	for key, value := range fileMap {
		exists := curEnvMap[key]
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
