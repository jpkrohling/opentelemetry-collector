package configgrpc

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"

	"google.golang.org/grpc/credentials"
)

var _ credentials.PerRPCCredentials = (*PerRPCAuth)(nil)

// PerRPCAuth is a gRPC credentials.PerRPCCredentials implementation that returns an 'authorization' header
type PerRPCAuth struct {
	metadata map[string]string
}

// BearerTokenFromFile builds a new PerRPCAuth with bearer token authentication, reading the token from the specified file
func BearerTokenFromFile(file string) (*PerRPCAuth, error) {
	token, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("couldn't obtain token from file: %w", err)
	}

	re := regexp.MustCompile(`[\s\n]`)
	token = re.ReplaceAll(token, []byte(""))

	return &PerRPCAuth{
		metadata: map[string]string{"authorization": fmt.Sprintf("Bearer %s", token)},
	}, nil
}

// BearerToken returns a new PerRPCAuth based on the given token
func BearerToken(t string) *PerRPCAuth {
	return &PerRPCAuth{
		metadata: map[string]string{"authorization": fmt.Sprintf("Bearer %s", t)},
	}
}

// GetRequestMetadata returns the request metadata to be used with the RPC
func (c *PerRPCAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return c.metadata, nil
}

// RequireTransportSecurity always returns true for this implementation. Passing bearer tokens in plain-text connections is a bad idea.
func (c *PerRPCAuth) RequireTransportSecurity() bool {
	return false
}
