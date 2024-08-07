package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	ipKey = "X-Real-IP"
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
			ipStr := r.Header.Get(ipKey)
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

// SubnetInterceptor is a gRPC interceptor that verifies if the IP address of the incoming
// request is within a trusted subnet.
func (s *Server) SubnetInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	needCheckTrustedSubnet := s.trustedSubnet != nil
	if needCheckTrustedSubnet {
		var ipStr string

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "subnet interceptor: missing request metadata")
		}

		values := md.Get(ipKey)
		if len(values) > 0 {
			ipStr = values[0]
		}
		if len(ipStr) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "subnet interceptor: missing ip address")
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, status.Errorf(codes.InvalidArgument, "subnet interceptor: invalid ip address")
		}

		if !s.trustedSubnet.Contains(ip) {
			return nil, status.Errorf(codes.Unauthenticated, "ip address isn`t in trusted subnet")
		}
	}

	return handler(ctx, req)
}
