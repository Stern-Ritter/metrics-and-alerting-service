package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Random struct {
	rand *rand.Rand
}

func NewRandom() Random {
	s := rand.NewSource(time.Now().Unix())
	return Random{rand.New(s)}
}

func (r *Random) Float(min, max float64) (float64, error) {
	if max < min {
		return 0, fmt.Errorf("min value: %f should be less than max value: %f", min, max)
	}

	res := min + r.rand.Float64()*(max-min)
	return res, nil
}

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

func CopyMap[K comparable, V any](source map[K]V) map[K]V {
	result := make(map[K]V)

	for key, value := range source {
		result[key] = value
	}

	return result
}

func FormatGaugeMetricValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func FormatCounterMetricValue(value int64) string {
	return strconv.Itoa(int(value))
}

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

func AddProtocolPrefix(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return strings.Join([]string{"http://", url}, "")
	}
	return url
}
