package goenv

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

var (
	missingDelimeter error = fmt.Errorf("Missing '=' in environment variable")
	missingEndQuote  error = fmt.Errorf("Missing end quote '\"' in environment variable")
)

func parseInput(src []byte) (map[string]string, error) {
	srcMap := make(map[string]string)
	src = bytes.ReplaceAll(src, []byte("\r\n"), []byte("\n"))

	for {
		rest := findLineStart(src)
		if rest == nil {
			// end of file
			break
		}

		key, rest, err := findKey(rest)
		if err != nil {
			return nil, err
		}

		value, rest, err := findValue(rest)
		if err != nil {
			return nil, err
		}

		srcMap[key] = value
		src = rest
	}

	return srcMap, nil
}

func findKey(src []byte) (string, []byte, error) {
	delimIndex := bytes.IndexFunc(src, func(r rune) bool {
		return r == '='
	})
	if delimIndex == -1 {
		return "", nil, missingDelimeter
	}

	key := string(src[:delimIndex])
	rest := src[delimIndex+1:]

	return strings.TrimSpace(key), rest, nil
}

func findValue(src []byte) (string, []byte, error) {
	src, isQuoted := findValueStart(src)
	if isQuoted {
		return readStringValue(src)
	}

	delimIndex := bytes.IndexFunc(src, func(r rune) bool {
		return r == '\n'
	})

	var value []byte
	var rest []byte
	if delimIndex == -1 {
		value = src
		rest = []byte{}
	} else {
		value = src[:delimIndex]
		rest = src[delimIndex+1:]
	}

	// check line for end of line comment
	valLength := len(value)
	// if first char of value == '#' then we allow it to be there for now
	for i := 1; i < valLength; i++ {
		if value[i] == '#' {
			if isSpace(value[i-1]) {
				valLength = i
				break
			}
		}
	}

	return strings.TrimSpace(string(value[:valLength])), rest, nil
}

func readStringValue(src []byte) (string, []byte, error) {
	endQuoteIndex := bytes.IndexFunc(src, func(r rune) bool {
		return r == '"'
	})
	if endQuoteIndex == -1 {
		return "", nil, missingEndQuote
	}

	value := string(src[:endQuoteIndex])
	rest := src[endQuoteIndex+1:]

	delimIndex := bytes.IndexFunc(rest, func(r rune) bool {
		return r == '\n'
	})
	if delimIndex != -1 {
		rest = rest[delimIndex+1:]
	}

	return value, rest, nil
}

func findValueStart(src []byte) ([]byte, bool) {
	nonSpaceIndex := bytes.IndexFunc(src, func(r rune) bool {
		return !unicode.IsSpace(r)
	})
	if nonSpaceIndex == -1 {
		return nil, false
	}
	src = src[nonSpaceIndex:]
	isQuoted := src[0] == '"'

	if isQuoted {
		src = src[1:]
	}

	return src, isQuoted
}

func findLineStart(src []byte) []byte {
	nonSpaceIndex := bytes.IndexFunc(src, func(r rune) bool {
		return !unicode.IsSpace(r)
	})
	if nonSpaceIndex == -1 {
		return nil
	}

	src = src[nonSpaceIndex:]
	if src[0] != '#' {
		return src
	}

	newLineIndex := bytes.IndexFunc(src, func(r rune) bool {
		return r == '\n'
	})
	if newLineIndex == -1 {
		return nil
	}

	return findLineStart(src[newLineIndex:])
}

// isSpace reports whether the rune is a space character but not line break character
//
// this differs from unicode.IsSpace, which also applies line break as space
func isSpace(r byte) bool {
	switch r {
	case '\t', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	}
	return false
}
