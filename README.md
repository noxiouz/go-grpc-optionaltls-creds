![workflow](https://github.com/noxiouz/go-grpc-optionaltls-creds/actions/workflows/go.yml/badge.svg)

# go-grpc-optionaltls-creds
The Go language implementation  of gRPC TransportCredentials that supports optional TLS connections.


# Simple example

**optionaltls.New** wraps provided TransportCredentials. It uses provided credentials if a client connected with TLS and bypasses if connection is plain-text

```golang
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    "github.com/noxiouz/go-grpc-optionaltls-creds/optionaltls"
)

// Control in runtime. See docs
var dynamicOptionF = optionaltls.DynamicOption(func() bool {
    return true
})

func createServer(credentials credentials.TransportCredentials) *grpc.Server {
    serverCredentials = optionaltls.New(credentials) // optionaltls.NewWithDynamicOption(credentials, dynamicOptionF))
    s := grpc.NewServer(grpc.Creds(serverCredentials))
    return s
}
```

# Implementation

Detection mechanism is inspired by [fbthrift](https://github.com/facebook/fbthrift/blob/master/thrift/lib/cpp2/server/peeking/TLSHelper.cpp#L29) and [The Illustrated TLS Connection](https://tls.ulfheim.net/)
