package nats

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nats-io/nats.go"
	"github.com/makesalekz/utils/v1/config"
)

type IQueue interface {
	Pub(data any)
}

type IQueueManager interface {
	GetLocal(name string) IQueue
	GetRemote(subj string) IQueue
	AddConsumer(name string, handler func(ctx context.Context, m *nats.Msg) bool)
	AddRemoteConsumer(name string, queueName string, handler func(ctx context.Context, m *nats.Msg) bool)
}

type QueueManager struct {
	nc      *nats.EncodedConn
	log     *log.Helper
	appName string
	ql      sync.RWMutex
	queues  map[string]*Queue
}

func NewQueueManager(c *config.Config, nc *nats.EncodedConn, logger log.Logger) IQueueManager {
	return &QueueManager{
		nc:      nc,
		log:     log.NewHelper(logger),
		appName: c.GetAppName(),
		queues:  make(map[string]*Queue),
	}
}

type queueKey struct{}

func (qm *QueueManager) GetLocal(name string) IQueue {
	subj := qm.appName + "/" + name

	queue := qm.getQueue(subj)
	if queue == nil {
		queue = qm.initQueue(subj)
	}

	return queue
}

func (qm *QueueManager) GetRemote(subj string) IQueue {
	queue := qm.getQueue(subj)
	if queue == nil {
		queue = qm.initQueue(subj)
	}

	return queue
}

func (qm *QueueManager) AddConsumer(name string, handler func(ctx context.Context, m *nats.Msg) bool) {
	subj := qm.appName + "/" + name
	queue := qm.GetLocal(name)

	ctx := context.WithValue(context.Background(), queueKey{}, queue)
	_, err := qm.nc.QueueSubscribe(subj, "workers", func(m *nats.Msg) {
		if !handler(ctx, m) {
			m.Nak()
		}
	})

	if err != nil {
		qm.log.Errorf("nc.QueueSubscribe: %s", err.Error())
	}
}

func (qm *QueueManager) AddRemoteConsumer(name string, queueName string, handler func(ctx context.Context, m *nats.Msg) bool) {
	queue := qm.GetRemote(name)

	ctx := context.WithValue(context.Background(), queueKey{}, queue)
	_, err := qm.nc.QueueSubscribe(name, queueName, func(m *nats.Msg) {
		if !handler(ctx, m) {
			m.Nak()
		}
	})

	if err != nil {
		qm.log.Errorf("nc.QueueSubscribe: %s", err.Error())
	}
}

func (qm *QueueManager) getQueue(subj string) *Queue {
	qm.ql.RLock()
	defer qm.ql.RUnlock()
	return qm.queues[subj]
}

func (qm *QueueManager) initQueue(subj string) *Queue {
	qm.ql.Lock()
	defer qm.ql.Unlock()

	queue, ok := qm.queues[subj]
	if !ok {
		queue = newQueue(qm.nc, qm.log, subj)
		qm.queues[subj] = queue
	}

	return queue
}

type Queue struct {
	nc   *nats.EncodedConn
	log  *log.Helper
	name string
}

func newQueue(nc *nats.EncodedConn, log *log.Helper, name string) *Queue {
	return &Queue{
		nc:   nc,
		log:  log,
		name: name,
	}
}

func (q *Queue) Pub(data any) {
	err := q.nc.Publish(q.name, data)
	if err != nil {
		q.log.Warnf("nc.Publish for %s: %s", q.name, err.Error())
	}
}
