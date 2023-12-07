package nats

import (
	"github.com/nats-io/nats.go"
	"gitlab.calendaria.team/services/utils/v1/config"
)

type NatsOptions func(o *NatsClient)

func WithUrl(url string) NatsOptions {
	return func(o *NatsClient) {
		o.natsUrl = url
	}
}

type NatsClient struct {
	*nats.EncodedConn
	natsUrl string
}

// NewNatsClient .
func NewNatsClient(c *config.Config, opts ...NatsOptions) (*NatsClient, func(), error) {
	natsClient := &NatsClient{}
	for _, opt := range opts {
		opt(natsClient)
	}

	nc, err := nats.Connect(natsClient.natsUrl)
	if err != nil {
		return nil, nil, err
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		nc.Close()
		return nil, nil, err
	}

	cleanup := func() {
		ec.Close()
		nc.Close()
	}

	natsClient.EncodedConn = ec

	return natsClient, cleanup, nil
}
