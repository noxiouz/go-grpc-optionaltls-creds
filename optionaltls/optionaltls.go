package optionaltls

import (
	"context"
	"net"

	"google.golang.org/grpc/credentials"
)

type optionalTLSCreds struct {
	tc credentials.TransportCredentials
}

func (c *optionalTLSCreds) Info() credentials.ProtocolInfo {
	return c.tc.Info()
}

func (c *optionalTLSCreds) Clone() credentials.TransportCredentials {
	return New(c.tc.Clone())
}

func (c *optionalTLSCreds) OverrideServerName(name string) error {
	return c.tc.OverrideServerName(name)
}

func (c *optionalTLSCreds) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	isTLS, bytes, err := DetectTLS(conn)
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	var wc net.Conn = NewWrappedConn(conn, bytes)
	if isTLS {
		return c.tc.ServerHandshake(wc)
	}

	var authInfo = info{
		CommonAuthInfo: credentials.CommonAuthInfo{SecurityLevel: credentials.NoSecurity},
	}

	return wc, authInfo, nil
}

func (c *optionalTLSCreds) ClientHandshake(ctx context.Context, authority string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return c.tc.ClientHandshake(ctx, authority, conn)
}

func New(tc credentials.TransportCredentials) credentials.TransportCredentials {
	return &optionalTLSCreds{tc: tc}
}

type info struct {
	credentials.CommonAuthInfo
}

func (info) AuthType() string {
	return "insecure"
}
