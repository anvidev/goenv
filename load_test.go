package goenv

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name          string
		shouldFail    bool
		files         []string
		expected      map[string]string
		expectedError error
	}{
		{
			name:  "Load default file",
			files: []string{},
			expected: map[string]string{
				"APP_ENV":             "production",
				"DEBUG":               "false",
				"API_URL":             "https://api.example.com",
				"DB_HOST":             "db.prod.internal",
				"DB_PORT":             "5432",
				"DB_USER":             "prod_user",
				"FEATURE_FLAG_NEW_UI": "false",
				"CACHE_TTL":           "300",
				"LOG_LEVEL":           "warn",
			},
			shouldFail: false,
		},
		{
			name:  "Load single file",
			files: []string{".env"},
			expected: map[string]string{
				"APP_ENV":             "production",
				"DEBUG":               "false",
				"API_URL":             "https://api.example.com",
				"DB_HOST":             "db.prod.internal",
				"DB_PORT":             "5432",
				"DB_USER":             "prod_user",
				"FEATURE_FLAG_NEW_UI": "false",
				"CACHE_TTL":           "300",
				"LOG_LEVEL":           "warn",
			},
			shouldFail: false,
		},
		{
			name:  "Load two files",
			files: []string{"testdata/.env.development", "testdata/.env.ci"},
			expected: map[string]string{
				"APP_ENV":             "development",
				"DEBUG":               "true",
				"API_URL":             "http://localhost:3000/api",
				"DB_HOST":             "localhost",
				"DB_PORT":             "5432",
				"DB_USER":             "dev_user",
				"FEATURE_FLAG_NEW_UI": "true",
				"CACHE_TTL":           "60",
				"LOG_LEVEL":           "debug",
				"RUN_E2E":             "true",
				"PARALLEL_JOBS":       "4",
				"GIT_COMMIT_SHA":      "abcdef123456",
				"CI_PIPELINE_ID":      "78910",
			},
			shouldFail: false,
		},
		{
			name:  "Load three files",
			files: []string{"testdata/.env.development", "testdata/.env.staging", "testdata/.env.ci"},
			expected: map[string]string{
				"APP_ENV":             "development",
				"DEBUG":               "true",
				"API_URL":             "http://localhost:3000/api",
				"DB_HOST":             "localhost",
				"DB_PORT":             "5432",
				"DB_USER":             "dev_user",
				"ENABLE_ANALYTICS":    "true",
				"ANALYTICS_KEY":       "stg-xyz-123",
				"MAINTENANCE_MODE":    "false",
				"FEATURE_FLAG_NEW_UI": "true",
				"CACHE_TTL":           "60",
				"LOG_LEVEL":           "debug",
				"RUN_E2E":             "true",
				"PARALLEL_JOBS":       "4",
				"GIT_COMMIT_SHA":      "abcdef123456",
				"CI_PIPELINE_ID":      "78910",
			},
			shouldFail: false,
		},
		{
			name:  "Load four files",
			files: []string{"testdata/.env.development", "testdata/.env.test", "testdata/.env.staging", "testdata/.env.ci"},
			expected: map[string]string{
				"APP_ENV":             "development",
				"DEBUG":               "true",
				"API_URL":             "http://localhost:3000/api",
				"DB_HOST":             "localhost",
				"DB_PORT":             "5432",
				"DB_USER":             "dev_user",
				"USE_MOCKS":           "true",
				"TIMEOUT_MS":          "5000",
				"RETRY_COUNT":         "3",
				"ENABLE_ANALYTICS":    "true",
				"ANALYTICS_KEY":       "stg-xyz-123",
				"MAINTENANCE_MODE":    "false",
				"FEATURE_FLAG_NEW_UI": "true",
				"CACHE_TTL":           "60",
				"LOG_LEVEL":           "debug",
				"RUN_E2E":             "true",
				"PARALLEL_JOBS":       "4",
				"GIT_COMMIT_SHA":      "abcdef123456",
				"CI_PIPELINE_ID":      "78910",
			},
			shouldFail: false,
		},
		{
			name:  "Load five files",
			files: []string{".env", "testdata/.env.development", "testdata/.env.test", "testdata/.env.staging", "testdata/.env.ci"},
			expected: map[string]string{
				"APP_ENV":             "production",
				"DEBUG":               "false",
				"API_URL":             "https://api.example.com",
				"DB_HOST":             "db.prod.internal",
				"DB_PORT":             "5432",
				"DB_USER":             "prod_user",
				"FEATURE_FLAG_NEW_UI": "false",
				"CACHE_TTL":           "300",
				"LOG_LEVEL":           "warn",
				"USE_MOCKS":           "true",
				"TIMEOUT_MS":          "5000",
				"RETRY_COUNT":         "3",
				"ENABLE_ANALYTICS":    "true",
				"ANALYTICS_KEY":       "stg-xyz-123",
				"MAINTENANCE_MODE":    "false",
				"RUN_E2E":             "true",
				"PARALLEL_JOBS":       "4",
				"GIT_COMMIT_SHA":      "abcdef123456",
				"CI_PIPELINE_ID":      "78910",
			},
			shouldFail: false,
		},
		{
			name:          "Load non existing file",
			files:         []string{"testdata/.env.missing"},
			shouldFail:    true,
			expectedError: fmt.Errorf("goenv: Failed to load file 'testdata/.env.missing': open testdata/.env.missing: no such file or directory"),
		},
		{
			name:          "Load file with missing '='",
			files:         []string{"testdata/.env.missing.delim"},
			shouldFail:    true,
			expectedError: fmt.Errorf("goenv: Failed to load file 'testdata/.env.missing.delim': malformed line: Missing '=' in environment variable"),
		},
		{
			name:          "Load file with malformed string",
			files:         []string{"testdata/.env.malformed"},
			shouldFail:    true,
			expectedError: fmt.Errorf("goenv: Failed to load file 'testdata/.env.malformed': malformed value: Missing end quote '\"' in environment variable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Load(tt.files...)
			if tt.shouldFail {
				if err == nil {
					t.Errorf("Load() = did not fail when expected.")
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("Load() = got error:'%s', but expected '%s'", err.Error(), tt.expectedError.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Load() = failed with error: %v", err)
				}
				got := make(map[string]string)

				for key := range tt.expected {
					got[key] = os.Getenv(key)
				}

				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Load() = got:\n%v\nexpected:\n%v", got, tt.expected)
				}
			}
		})
	}
}
