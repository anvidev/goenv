package goenv

import (
	"bytes"
	"io"
	"os"
)

func Load(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = append(filenames, ".env")
	}

	for _, filename := range filenames {
		if err := loadFile(filename); err != nil {
			return err
		}
	}

	return nil
}

func loadFile(filename string) error {
	src, err := readFile(filename)
	if err != nil {
		return err
	}

	fileMap, err := parseInput(src)
	if err != nil {
		return err
	}

	for key, value := range fileMap {
		os.Setenv(key, value)
	}

	return nil
}

func readFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, f)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
