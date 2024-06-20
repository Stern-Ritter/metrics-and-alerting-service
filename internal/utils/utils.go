package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Random generates random numbers
type Random struct {
	rand *rand.Rand
}

// NewRandom is constructor for creating an instance of Random
func NewRandom() Random {
	s := rand.NewSource(time.Now().Unix())
	return Random{rand.New(s)}
}

// Float generates a random float64 number in the specified range [min, max].
// It returns an error if max is less than min.
func (r *Random) Float(min, max float64) (float64, error) {
	if max < min {
		return 0, fmt.Errorf("random generator min value: %f should be less than max value: %f", min, max)
	}

	res := min + r.rand.Float64()*(max-min)
	return res, nil
}

// Contains checks if the slice s contains any of the values v.
// The comparison is case-insensitive.
func Contains(s []string, v ...string) bool {
	for _, e := range s {
		for _, f := range v {
			if strings.EqualFold(e, f) {
				return true
			}
		}
	}
	return false
}

// CopyMap creates a shallow copy of the map.
// It returns a new map with the same key-value pairs as the source map.
func CopyMap[K comparable, V any](source map[K]V) map[K]V {
	result := make(map[K]V, len(source))

	for key, value := range source {
		result[key] = value
	}

	return result
}

// FormatGaugeMetricValue formats a float64 gauge metric value as a string.
func FormatGaugeMetricValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// FormatCounterMetricValue formats an int64 counter metric value as a string.
func FormatCounterMetricValue(value int64) string {
	return strconv.Itoa(int(value))
}

// ValidateHostnamePort checks if the given string is in the format "hostname:port".
// It returns an error if the format is invalid.
func ValidateHostnamePort(hp string) error {
	pattern := `[^\:]+:[0-9]{1,5}`
	regexp, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	if !regexp.MatchString(hp) {
		return fmt.Errorf("invalid hostname and port format:%s, hould be: <host>:<port>", hp)
	}
	return nil
}

// AddProtocolPrefix adds prefix "http://" to the URL if it does not already start with "http://" or "https://".
// It returns the modified URL.
func AddProtocolPrefix(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return strings.Join([]string{"http://", url}, "")
	}
	return url
}
