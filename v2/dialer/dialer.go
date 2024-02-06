package dialer

import (
	"context"
	"time"

	"gitlab.calendaria.team/services/utils/v2/middlewares/auth"

	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const CONNECTING_TIMEOUT = 3 * time.Second
const TRANSIENT_FAILURE_TIMEOUT = 5 * time.Second

type IDialer interface {
	SetEndpoint(endpointName, endpoint string) IDialer
	SetTimeout(timeout time.Duration) IDialer
	Connect(ctx context.Context) (*ggrpc.ClientConn, error)
	Close() error
}

// Dialer is a service dialer, implements IDialer
type Dialer struct {
	conn        *ggrpc.ClientConn
	dm          *DialerManager
	jwtAudience jwtv4.ClaimStrings
	endpoint    string
	timeout     time.Duration
}

func (d *Dialer) SetEndpoint(endpointName, endpoint string) IDialer {
	d.endpoint = endpoint
	d.jwtAudience = jwtv4.ClaimStrings{endpointName}

	return d
}

func (d *Dialer) SetTimeout(timeout time.Duration) IDialer {
	d.timeout = timeout

	return d
}

func (d *Dialer) Connect(ctx context.Context) (*ggrpc.ClientConn, error) {
	if d.conn != nil {
		s := d.conn.GetState()
		switch s {
		case connectivity.Idle:
			// we can use this connection
			return d.conn, nil
		case connectivity.Ready:
			// we can use this connection
			return d.conn, nil
		case connectivity.Connecting:
			// we should wait for connection
			waitCtx, cancelWait := context.WithTimeout(ctx, CONNECTING_TIMEOUT)
			defer cancelWait()
			if d.conn.WaitForStateChange(waitCtx, s) {
				return d.conn, nil
			}
		case connectivity.TransientFailure:
			// we should wait for connection
			waitCtx, cancelWait := context.WithTimeout(ctx, TRANSIENT_FAILURE_TIMEOUT)
			defer cancelWait()
			if d.conn.WaitForStateChange(waitCtx, s) {
				return d.conn, nil
			}
		case connectivity.Shutdown:
			// we should reconnect
		}
	}

	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(d.endpoint),
		grpc.WithDiscovery(d.dm.discovery),
		grpc.WithTimeout(d.timeout),
		grpc.WithMiddleware(
			auth.Client(d.dm.jwt, d.dm.jwtIssuer, d.jwtAudience),
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
