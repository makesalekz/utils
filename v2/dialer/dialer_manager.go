package dialer

import (
	"time"

	"gitlab.calendaria.team/services/utils/v1/config"
	jwtp "gitlab.calendaria.team/services/utils/v1/jwt"

	consul "github.com/go-kratos/consul/registry"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

type IDialerManager interface {
	NewServiceDialer(endpointName string, endpoint string) (IDialer, error)
}

// DialerManager is a service dialer manager
type DialerManager struct {
	discovery *consul.Registry
	jwt       *jwtp.JwtProcessor
	jwtIssuer string
}

func NewServiceDialerManager(
	c *config.Config,
	jwt *jwtp.JwtProcessor,
) (IDialerManager, error) {
	return &DialerManager{
		discovery: c.GetRegistry(),
		jwt:       jwt,
		jwtIssuer: c.GetAppName(),
	}, nil
}

func (dm *DialerManager) NewServiceDialer(
	endpointName string,
	endpoint string,
) (IDialer, error) {
	return &Dialer{
		dm:          dm,
		jwtAudience: jwtv4.ClaimStrings{endpointName},
		endpoint:    endpoint,
		timeout:     30 * time.Second,
	}, nil
}
