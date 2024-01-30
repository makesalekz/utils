package data

import (
	"context"
	"fmt"
	"os"
	"time"

	"gitlab.calendaria.team/services/utils/v1/config"
	jwtp "gitlab.calendaria.team/services/utils/v1/jwt"
	"gitlab.calendaria.team/services/utils/v2/middlewares/auth"

	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	ggrpc "google.golang.org/grpc"
)

type Dialer struct {
	conn        *ggrpc.ClientConn
	discovery   *consul.Registry
	jwt         *jwtp.JwtProcessor
	jwtIssuer   string
	jwtAudience jwtv4.ClaimStrings
	endpoint    string
	timeout     time.Duration
}

func NewServiceDialer(
	c *config.Config,
	jwt *jwtp.JwtProcessor,
	endpointName string,
	endpoint string,
) (*Dialer, error) {
	appName := os.Getenv("SERVICE_NAME")
	if appName == "" {
		return nil, fmt.Errorf("SERVICE_NAME not found")
	}

	return &Dialer{
		discovery:   c.GetRegistry(),
		jwt:         jwt,
		jwtIssuer:   appName,
		jwtAudience: jwtv4.ClaimStrings{endpointName},
		endpoint:    endpoint,
		timeout:     30 * time.Second,
	}, nil
}

func (d *Dialer) SetEndpoint(endpointName, endpoint string) *Dialer {
	d.endpoint = endpoint
	d.jwtAudience = jwtv4.ClaimStrings{endpointName}

	return d
}

func (d *Dialer) SetTimeout(timeout time.Duration) *Dialer {
	d.timeout = timeout

	return d
}

func (d *Dialer) Connect(ctx context.Context) (*ggrpc.ClientConn, error) {
	if d.conn != nil {
		return d.conn, nil
	}

	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(d.endpoint),
		grpc.WithDiscovery(d.discovery),
		grpc.WithTimeout(d.timeout),
		grpc.WithMiddleware(
			auth.Client(d.jwt, d.jwtIssuer, d.jwtAudience),
			metadata.Client(),
		),
	)
	if err != nil {
		return nil, err
	}

	d.conn = conn

	return conn, nil
}

func (d *Dialer) Close() error {
	if d.conn == nil {
		return nil
	}

	return d.conn.Close()
}
