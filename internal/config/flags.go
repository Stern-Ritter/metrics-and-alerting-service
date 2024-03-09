package config

import (
	"fmt"
	"strconv"
	"strings"
)

type URL struct {
	Host string
	Port int
}

func (u *URL) String() string {
	return fmt.Sprintf("%s:%s", u.Host, strconv.Itoa(u.Port))
}

func (u *URL) Set(flagValue string) error {
	hp := strings.Split(flagValue, ":")

	if len(hp) != 2 {
		return fmt.Errorf("invalid url value error: format should be <host>:<port>")
	}

	host := hp[0]
	port, err := strconv.Atoi(hp[1])

	if err != nil {
		return fmt.Errorf("invalid url port value error: should be a number")
	}

	u.Host = host
	u.Port = port

	return nil
}
