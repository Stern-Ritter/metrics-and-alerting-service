package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestGetTrustedSubnet(t *testing.T) {
	testCases := []struct {
		name           string
		cidrStr        string
		expectedSubnet *net.IPNet
		expectedErr    error
	}{
		{
			name:    "valid CIDR notated trusted subnet",
			cidrStr: "192.168.1.0/24",
			expectedSubnet: &net.IPNet{
				IP:   net.ParseIP("192.168.1.0"),
				Mask: net.CIDRMask(24, 32),
			},
			expectedErr: nil,
		},
		{
			name:           "empty CIDR notated trusted subnet",
			cidrStr:        "",
			expectedSubnet: nil,
			expectedErr:    nil,
		},
		{
			name:           "invalid CIDR notated trusted subnet",
			cidrStr:        "192.168.1.0.0",
			expectedSubnet: nil,
			expectedErr:    fmt.Errorf("parse subnet from CIDR"),
		},
		{
			name:           "CIDR notated trusted subnet with missing mask",
			cidrStr:        "192.168.1.0",
			expectedSubnet: nil,
			expectedErr:    fmt.Errorf("parse subnet from CIDR"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			subnet, err := GetTrustedSubnet(tt.cidrStr)

			if tt.expectedErr != nil {
				assert.Error(t, err, "expected error but got none")
				assert.Equal(t, tt.expectedErr, err, "expected error: %s, got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "expected no error but got: %s", err)
				if tt.expectedSubnet != nil {
					assert.NotNil(t, subnet, "expected subnet but got nil")
					assert.Equal(t, tt.expectedSubnet.String(), subnet.String(), "expected subnet: %s, got: %s", tt.expectedSubnet, subnet)
				} else {
					assert.Nil(t, subnet, "expected nil subnet but got: %v", subnet)
				}
			}
		})
	}
}

func TestSubnetMiddleware(t *testing.T) {
	trustedSubnet := &net.IPNet{
		IP:   net.ParseIP("192.168.1.0"),
		Mask: net.CIDRMask(24, 32),
	}

	tests := []struct {
		name           string
		trustedSubnet  *net.IPNet
		ipHeader       string
		expectedStatus int
	}{
		{
			name:           "IP is in trusted subnet",
			trustedSubnet:  trustedSubnet,
			ipHeader:       "192.168.1.10",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid IP address",
			trustedSubnet:  trustedSubnet,
			ipHeader:       "192.168.1.0.0",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "IP is outside trusted subnet",
			trustedSubnet:  trustedSubnet,
			ipHeader:       "192.168.2.10",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "No trusted subnet",
			trustedSubnet:  nil,
			ipHeader:       "192.168.1.10",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{trustedSubnet: tt.trustedSubnet}
			middleware := server.SubnetMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(ipKey, tt.ipHeader)
			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code, "expected http status code: %d, got: %d", tt.expectedStatus, rec.Code)
		})
	}
}

func TestSubnetInterceptor(t *testing.T) {
	trustedSubnet := &net.IPNet{
		IP:   net.ParseIP("192.168.1.0"),
		Mask: net.CIDRMask(24, 32),
	}

	tests := []struct {
		name          string
		trustedSubnet *net.IPNet
		md            metadata.MD
		expectedErr   error
	}{
		{
			name:          "IP in trusted subnet",
			trustedSubnet: trustedSubnet,
			md:            metadata.Pairs(ipKey, "192.168.1.10"),
			expectedErr:   nil,
		},
		{
			name:          "No metadata",
			trustedSubnet: trustedSubnet,
			md:            nil,
			expectedErr:   status.Errorf(codes.InvalidArgument, "subnet interceptor: missing request metadata"),
		},
		{
			name:          "Invalid IP address",
			trustedSubnet: trustedSubnet,
			md:            metadata.Pairs(ipKey, "192.168.1.0.0"),
			expectedErr:   status.Errorf(codes.InvalidArgument, "subnet interceptor: invalid ip address"),
		},
		{
			name:          "Missing IP address",
			trustedSubnet: trustedSubnet,
			md:            metadata.Pairs(),
			expectedErr:   status.Errorf(codes.InvalidArgument, "subnet interceptor: missing ip address"),
		},
		{
			name:          "IP outside trusted subnet",
			trustedSubnet: trustedSubnet,
			md:            metadata.Pairs(ipKey, "192.168.2.10"),
			expectedErr:   status.Errorf(codes.Unauthenticated, "ip address isn`t in trusted subnet"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{trustedSubnet: tt.trustedSubnet}

			ctx := context.Background()
			if tt.md != nil {
				ctx = metadata.NewIncomingContext(ctx, tt.md)
			}
			_, err := server.SubnetInterceptor(ctx, nil, nil,
				func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil })

			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err, "should return error: %s, got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "should not return error, but got: %s", err)
			}
		})
	}
}
