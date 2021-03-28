package tests

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/noxiouz/go-grpc-optionaltls-creds/optionaltls"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func createUnstartedServer(creds credentials.TransportCredentials) *grpc.Server {
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterGreeterServer(s, &server{})
	return s
}

type testCredentials struct {
	client credentials.TransportCredentials
	server credentials.TransportCredentials
}

func createCredentials() (*testCredentials, error) {
	cert, err := tls.X509KeyPair(localhostCert, localhostKey)
	if err != nil {
		return nil, err
	}

	certificate, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	certpool := x509.NewCertPool()
	certpool.AddCert(certificate)

	tc := &testCredentials{
		client: credentials.NewClientTLSFromCert(certpool, "example.com"),
		server: credentials.NewServerTLSFromCert(&cert),
	}
	return tc, nil
}

func TestOptionalTLS(t *testing.T) {
	testCtx, testCancel := context.WithCancel(context.Background())
	defer testCancel()

	tc, err := createCredentials()
	if err != nil {
		t.Fatalf("failed to create credentials %v", err)
	}

	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("failed to listen %v", err)
	}
	defer lis.Close()
	addr := lis.Addr().String()

	srv := createUnstartedServer(optionaltls.New(tc.server))
	go func() {
		srv.Serve(lis)
	}()
	defer srv.Stop()

	testFunc := func(t *testing.T, dialOpt grpc.DialOption) {
		ctx, cancel := context.WithTimeout(testCtx, 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, addr, dialOpt)
		if err != nil {
			t.Fatalf("failed to connect to the server %v", err)
		}
		defer conn.Close()
		c := pb.NewGreeterClient(conn)
		resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: "noxiouz"})
		if err != nil {
			t.Fatalf("could not greet: %v", err)
		}
		if resp.Message != "Hello noxiouz" {
			t.Fatalf("unexpected reply %s", resp.Message)
		}
	}

	t.Run("Plain2TLS", func(t *testing.T) {
		testFunc(t, grpc.WithInsecure())
	})
	t.Run("TLS2TLS", func(t *testing.T) {
		testFunc(t, grpc.WithTransportCredentials(tc.client))
	})
}
