package goenv

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expected      map[string]string
		shouldFail    bool
		expectedError error
	}{
		{
			name: "Parse nicely formatted input",
			input: []byte(` 
			ENVIRONMENT=development
			HOST=localhost
			`),
			expected: map[string]string{
				"ENVIRONMENT": "development",
				"HOST":        "localhost",
			},
		},
		{
			name: "Parse input with many spaces",
			input: []byte(` 
			  ENVIRONMENT  =   development
					HOST	=  localhost	
			`),
			expected: map[string]string{
				"ENVIRONMENT": "development",
				"HOST":        "localhost",
			},
		},
		{
			name: "Parse input with line comment",
			input: []byte(` 
			ENVIRONMENT=development
			# This is a comment
		  HOST=localhost	
			`),
			expected: map[string]string{
				"ENVIRONMENT": "development",
				"HOST":        "localhost",
			},
		},
		{
			name: "Parse input with comment at end of line",
			input: []byte(` 
			ENVIRONMENT= development #This is a comment
		  HOST=        localhost	
			`),
			expected: map[string]string{
				"ENVIRONMENT": "development",
				"HOST":        "localhost",
			},
		},
		{
			name: "Parse input with no trailing newline",
			input: []byte(` 
			ENVIRONMENT=development
		  HOST=localhost`),
			expected: map[string]string{
				"ENVIRONMENT": "development",
				"HOST":        "localhost",
			},
		},
		{
			name: "Parse input in quotes",
			input: []byte(` 
			ENVIRONMENT=  "development"
			API_BASE_URL= "https://example.com/api"
			`),
			expected: map[string]string{
				"ENVIRONMENT":  "development",
				"API_BASE_URL": "https://example.com/api",
			},
		},
		{
			name: "Parse input in quotes with inline comments",
			input: []byte(` 
			ENVIRONMENT=  " development "            # Comment with lots of spaces
			API_BASE_URL= "https://example.com/api"  # Another comment
			`),
			expected: map[string]string{
				"ENVIRONMENT":  " development ",
				"API_BASE_URL": "https://example.com/api",
			},
		},
		{
			name: "Parse input with comments inside quotes",
			input: []byte(` 
			ENVIRONMENT=  "development # a 'comment' inside a quoted value"
			API_BASE_URL= "https://example.com/api"  # Another comment
			`),
			expected: map[string]string{
				"ENVIRONMENT":  "development # a 'comment' inside a quoted value",
				"API_BASE_URL": "https://example.com/api",
			},
		},
		{
			name: "Parse input empty lines",
			input: []byte(` 

			ENVIRONMENT="development"

			API_BASE_URL="https://example.com/api"

			`),
			expected: map[string]string{
				"ENVIRONMENT":  "development",
				"API_BASE_URL": "https://example.com/api",
			},
		},
		{
			name: "Parse input ending with empty line",
			input: []byte(` 
			ENVIRONMENT="development"
			API_BASE_URL="https://example.com/api"

			`),
			expected: map[string]string{
				"ENVIRONMENT":  "development",
				"API_BASE_URL": "https://example.com/api",
			},
		},
		{
			name: "Parse input with missing delimeter",
			input: []byte(` 
			ENVIRONMENT="development"
			API_BASE_URL:"https://example.com/api"
			`),
			shouldFail:    true,
			expectedError: fmt.Errorf("malformed line: Missing '=' in environment variable"),
		},
		{
			name: "Parse input with malformed string value",
			input: []byte(` 
			ENVIRONMENT="development"
			API_BASE_URL="https://example.com/api
			`),
			shouldFail:    true,
			expectedError: fmt.Errorf("malformed value: Missing end quote '\"' in environment variable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInput(tt.input)
			if tt.shouldFail {
				if err == nil {
					t.Errorf("parseInput() = did not fail when expected.")
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("parseInput() = got error:'%s', but expected '%s'", err.Error(), tt.expectedError.Error())
				}
			} else {
				if err != nil {
					t.Errorf("parseInput() = failed with error: %v", err)
				}

				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("parseInput() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}
