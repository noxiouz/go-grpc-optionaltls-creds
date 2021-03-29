package optionaltls

import (
	"context"
	"net"

	"google.golang.org/grpc/credentials"
)

type optionalTLSCreds struct {
	tc            credentials.TransportCredentials
	dynamicOption DynamicOption
}

func (c *optionalTLSCreds) Info() credentials.ProtocolInfo {
	return c.tc.Info()
}

func (c *optionalTLSCreds) Clone() credentials.TransportCredentials {
	return NewWithDynamicOption(c.tc.Clone(), c.dynamicOption)
}

func (c *optionalTLSCreds) OverrideServerName(name string) error {
	return c.tc.OverrideServerName(name)
}

func (c *optionalTLSCreds) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	if c.dynamicOption != nil && !c.dynamicOption.IsActive() {
		return c.tc.ServerHandshake(conn)
	}

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
	return NewWithDynamicOption(tc, nil)
}

func NewWithDynamicOption(tc credentials.TransportCredentials, do DynamicOption) credentials.TransportCredentials {
	return &optionalTLSCreds{tc: tc, dynamicOption: do}
}

type info struct {
	credentials.CommonAuthInfo
}

func (info) AuthType() string {
	return "insecure"
}
