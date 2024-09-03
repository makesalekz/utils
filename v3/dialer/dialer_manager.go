package dialer

import (
	"time"

	u_config "gitlab.calendaria.team/services/utils/v1/config"
	u_tracing "gitlab.calendaria.team/services/utils/v2/tracing"
)

type IDialerManager interface {
	NewServiceDialer(endpointName string, endpoint string) (IDialer, error)
}

// DialerManager is a service dialer manager
type DialerManager struct {
	config *u_config.Config
	tracer *u_tracing.Tracer
}

func NewServiceDialerManager(
	config *u_config.Config,
	tracer *u_tracing.Tracer,
) (IDialerManager, error) {
	return &DialerManager{
		config: config,
		tracer: tracer,
	}, nil
}

func (dm *DialerManager) NewServiceDialer(
	endpointName string,
	endpoint string,
) (IDialer, error) {
	return &Dialer{
		dm:           dm,
		endpointName: endpointName,
		endpoint:     endpoint,
		timeout:      30 * time.Second,
	}, nil
}
