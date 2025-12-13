package client

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"os"

	client "github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/interseptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//go:embed certs/ca-cert.pem
var caCertPEM []byte

type GRPCConnection interface {
	Close() error
	Conn() *grpc.ClientConn
}

type grpcConnection struct {
	conn *grpc.ClientConn
}

func NewGRPCConnection(serverAddr string, tokenProvider client.TokenProvider) (GRPCConnection, error) {
	creds, err := loadTLSCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
	}

	ai := interseptors.NewAuthInterceptor(tokenProvider)

	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(ai.Interceptor()),
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	return &grpcConnection{
		conn: conn,
	}, nil
}

func (c *grpcConnection) Close() error {
	return c.conn.Close()
}

func (c *grpcConnection) Conn() *grpc.ClientConn {
	return c.conn
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(caCertPEM) {
		data, err := os.ReadFile("certs/ca-cert.pem")
		if err == nil && certPool.AppendCertsFromPEM(data) {

		} else {
			return nil, fmt.Errorf("no valid CA certificate found")
		}
	}

	config := &tls.Config{
		RootCAs:    certPool,
		ServerName: "localhost",
		MinVersion: tls.VersionTLS12,
	}

	return credentials.NewTLS(config), nil
}
