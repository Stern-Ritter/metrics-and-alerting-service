package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandom(t *testing.T) {
	t.Run("should return error if max value is less than min value", func(t *testing.T) {
		r := NewRandom()
		incorrectExpectedMin := 99.99
		incorrectExpectedMax := 0.1

		_, err := r.Float(incorrectExpectedMin, incorrectExpectedMax)
		require.Error(t, err)
	})

	t.Run("should return random value between min and max values", func(t *testing.T) {
		r := NewRandom()
		expectedMin := 0.1
		expectedMax := 99.99
		iterationCount := 1000

		for i := 0; i < iterationCount; i++ {
			got, err := r.Float(expectedMin, expectedMax)
			require.NoError(t, err)
			require.GreaterOrEqual(t, got, expectedMin, "random value %f should be greater or equal than min value: %f", got, expectedMin)
			require.LessOrEqual(t, got, expectedMax, "random value %f should be less or equal than max value: %f", got, expectedMax)
		}
	})
}

func TestContains(t *testing.T) {
	testCases := []struct {
		name   string
		source []string
		value  string
		want   bool
	}{
		{
			name:   "return true when source contains value and case of the letters matches",
			source: []string{"first", "second"},
			value:  "first",
			want:   true,
		},
		{
			name:   "return true when source contains value but case of the letters doesn`t matches",
			source: []string{"first", "second"},
			value:  "First",
			want:   true,
		},
		{
			name:   "return false when source doesn`t contain value",
			source: []string{"first", "second"},
			value:  "third",
			want:   false,
		},
		{
			name:   "return false when source is empty",
			source: []string{},
			value:  "first",
			want:   false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := Contains(tt.source, tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCopyMap(t *testing.T) {
	testCases := []struct {
		name   string
		source map[string]string
		want   map[string]string
	}{
		{
			name:   "return correct copy when source map contains values #1",
			source: map[string]string{"1": "first", "2": "second"},
			want:   map[string]string{"1": "first", "2": "second"},
		},
		{
			name:   "return correct copy when source map contains values #2",
			source: map[string]string{"1": "first", "2": "second"},
			want:   map[string]string{"2": "second", "1": "first"},
		},
		{
			name:   "return empty map when source map is empty",
			source: make(map[string]string),
			want:   make(map[string]string),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := CopyMap(tt.source)
			assert.True(t, reflect.DeepEqual(tt.source, got))
		})
	}
}

func TestValidateHostnamePort(t *testing.T) {
	testCases := []struct {
		name      string
		hp        string
		wantError bool
	}{
		{
			name:      "return nil when valid hostname and port #1",
			hp:        "localhost:8080",
			wantError: false,
		},
		{
			name:      "return nil when valid hostname and port #2",
			hp:        "https://ya.ru:443",
			wantError: false,
		},
		{
			name:      "return error when hostname without port",
			hp:        "localhost",
			wantError: true,
		},
		{
			name:      "return error when port without hostname",
			hp:        ":8082",
			wantError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHostnamePort(tt.hp)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddProtocolPrefix(t *testing.T) {
	testCases := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "return url with prefix when ulr doesn`t have prefix",
			url:  "pkg.go.dev/",
			want: "http://pkg.go.dev/",
		},
		{
			name: "return the same url when ulr has prefix",
			url:  "http://pkg.go.dev/",
			want: "http://pkg.go.dev/",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := AddProtocolPrefix(tt.url)
			assert.Equal(t, tt.want, got)
		})
	}
}
