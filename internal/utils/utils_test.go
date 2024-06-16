package utils

import (
	"reflect"
	"strings"
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

func ContainsMap(s []string, v ...string) bool {
	m := make(map[string]struct{}, len(s))
	for _, e := range s {
		m[strings.ToLower(e)] = struct{}{}
	}
	for _, f := range v {
		if _, found := m[strings.ToLower(f)]; found {
			return true
		}
	}
	return false
}

func CopyMapWithoutInitSize[K comparable, V any](source map[K]V) map[K]V {
	result := make(map[K]V)
	for key, value := range source {
		result[key] = value
	}
	return result
}

func AddProtocolPrefixWithSb(url string) string {
	var sb strings.Builder
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		sb.WriteString("http://")
		sb.WriteString(url)
		return sb.String()
	}
	return url
}

func BenchmarkContains(b *testing.B) {
	s := []string{"a", "b", "c", "d", "e", "f", "g"}
	v := []string{"b", "d", "k"}

	b.Run("with nested loops", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Contains(s, v...)
		}
	})

	b.Run("with map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ContainsMap(s, v...)
		}
	})
}

func BenchmarkCopyMap(b *testing.B) {
	source := map[int]string{
		1: "a",
		2: "b",
		3: "c",
		4: "d",
		5: "e",
		6: "f",
		7: "g",
	}

	b.Run("with init size", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CopyMap(source)
		}
	})

	b.Run("without init size", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CopyMapWithoutInitSize(source)
		}
	})
}

func BenchmarkAddProtocolPrefix(b *testing.B) {
	url := "example.com"

	b.Run("with strings join", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AddProtocolPrefix(url)
		}
	})

	b.Run("with strings builder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AddProtocolPrefixWithSb(url)
		}
	})
}
