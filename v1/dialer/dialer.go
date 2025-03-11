package dialer

import (
	"context"
	"time"

	"gitlab.calendaria.team/services/utils/v1/config"
	jwtp "gitlab.calendaria.team/services/utils/v1/jwt"
	"gitlab.calendaria.team/services/utils/v1/middlewares/auth"

	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
)

type Dialer struct {
	conn      *ggrpc.ClientConn
	discovery *consul.Registry
	jwt       *jwtp.JwtProcessor
}

func NewDialer(c config.IConfig, jwt *jwtp.JwtProcessor) (*Dialer, error) {
	return &Dialer{
		discovery: c.GetRegistry(),
		jwt:       jwt,
	}, nil
}

func (d *Dialer) Close() error {
	if d.conn == nil {
		return nil
	}

	return d.conn.Close()
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

func (d *DialerBuilder[T]) Conn(ctx context.Context, defaultClaims *jwtp.TenantClaims) (T, error) {
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(d.endpoint),
		grpc.WithDiscovery(d.dialer.discovery),
		grpc.WithTimeout(d.timeout),
		grpc.WithMiddleware(auth.Client(ctx, d.dialer.jwt, defaultClaims)),
	)

	var nilVar T
	if err != nil {
		return nilVar, err
	}

	d.dialer.conn = conn

	return d.clientConn(conn), nil
}
