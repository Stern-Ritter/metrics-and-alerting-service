package server

import (
	"fmt"
	"net"
	"net/http"
)

// GetTrustedSubnet parses a CIDR notated trusted subnet and returns a parsed subnet.
// It returns an error if the CIDR notated trusted subnet is invalid.
func GetTrustedSubnet(cidrStr string) (*net.IPNet, error) {
	var trustedSubnet *net.IPNet

	needCheckTrustedSubnet := len(cidrStr) > 0
	if needCheckTrustedSubnet {
		_, subnet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return nil, fmt.Errorf("parse subnet from CIDR")
		}
		trustedSubnet = subnet
	}

	return trustedSubnet, nil
}

// SubnetMiddleware creates an HTTP middleware that checks if the request's IP
// address is within the trusted subnet. If the IP is not within the trusted
// subnet, the middleware responds with a 403 Forbidden status.
func (s *Server) SubnetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		needCheckTrustedSubnet := s.trustedSubnet != nil
		if needCheckTrustedSubnet {
			ipStr := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(ipStr)
			if ip == nil {
				http.Error(w, "Invalid ip address in X-Real-IP header", http.StatusForbidden)
				return
			}

			if !s.trustedSubnet.Contains(ip) {
				http.Error(w, "Ip address isn`t in trusted subnet", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
