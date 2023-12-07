package dialer

import (
	"context"
	"time"

	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/middleware"
	kjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"gitlab.calendaria.team/services/utils/v1/config"
	jwtp "gitlab.calendaria.team/services/utils/v1/jwt"
	ggrpc "google.golang.org/grpc"
)

type Dialer struct {
	discovery *consul.Registry
	jwt       *jwtp.JwtProcessor
}

// NewJwtProcessor .
func NewDialer(c *config.Config, jwt *jwtp.JwtProcessor) (*Dialer, error) {
	return &Dialer{
		discovery: c.GetRegistry(),
		jwt:       jwt,
	}, nil
}

func (d *Dialer) getJwtMiddleware(ctx context.Context) middleware.Middleware {
	return kjwt.Client(func(token *jwtv4.Token) (interface{}, error) {
		return d.jwt.GetSecret(), nil
	}, kjwt.WithSigningMethod(jwtv4.SigningMethodHS256), kjwt.WithClaims(func() jwtv4.Claims {
		claims, _ := d.jwt.GetClaimsFromContext(ctx)
		return claims
	}))
}

type DialerBuilder[T any] struct {
	dialer     *Dialer
	clientConn func(cc ggrpc.ClientConnInterface) T
	endpoint   string
	timeout    time.Duration
}

func NewDialerBuilder[T any](
	d *Dialer,
	clientCon func(cc ggrpc.ClientConnInterface) T,
) *DialerBuilder[T] {
	return &DialerBuilder[T]{
		clientConn: clientCon,
		dialer:     d,
	}
}

func (d *DialerBuilder[T]) SetEndpoint(endpoint string) *DialerBuilder[T] {
	d.endpoint = endpoint

	return d
}

func (d *DialerBuilder[T]) SetTimeout(timeout time.Duration) *DialerBuilder[T] {
	d.timeout = timeout

	return d
}

func (d *DialerBuilder[T]) Conn(ctx context.Context) (T, error) {
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(d.endpoint),
		grpc.WithDiscovery(d.dialer.discovery),
		grpc.WithTimeout(d.timeout),
		grpc.WithMiddleware(d.dialer.getJwtMiddleware(ctx)),
	)

	var nilVar T
	if err != nil {
		return nilVar, err
	}

	return d.clientConn(conn), nil
}
