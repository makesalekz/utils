package nats

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nats-io/nats.go"
	"gitlab.calendaria.team/services/utils/v1/config"
)

type QueueManager struct {
	nc      *nats.EncodedConn
	log     *log.Helper
	appName string
	ql      sync.RWMutex
	queues  map[string]*Queue
}

func NewQueueManager(c *config.Config, nc *nats.EncodedConn, logger log.Logger) *QueueManager {
	return &QueueManager{
		nc:      nc,
		log:     log.NewHelper(logger),
		appName: c.GetAppName(),
		queues:  make(map[string]*Queue),
	}
}

type queueKey struct{}

func (qm *QueueManager) GetLocal(name string) *Queue {
	subj := qm.appName + "/" + name

	queue := qm.getQueue(subj)
	if queue == nil {
		queue = qm.initQueue(subj)
	}

	return queue
}

func (qm *QueueManager) AddConsumer(name string, handler func(ctx context.Context, m *nats.Msg) bool) {
	subj := qm.appName + "/" + name
	queue := qm.GetLocal(name)

	//delays := []time.Duration{
	//	1 * time.Minute,
	//	5 * time.Minute,
	//	30 * time.Minute,
	//	3 * time.Hour,
	//}

	delays := []time.Duration{
		1 * time.Second,
		5 * time.Second,
		30 * time.Second,
		1 * time.Minute,
	}

	ctx := context.WithValue(context.Background(), queueKey{}, queue)
	_, err := qm.nc.QueueSubscribe(subj, "workers", func(m *nats.Msg) {

		if !handler(ctx, m) {
			retryCountStr := m.Header.Get("X-Retry-Count")
			if retryCountStr == "" {
				retryCountStr = "0"
			}
			retryCount, err := strconv.ParseInt(retryCountStr, 10, 64)
			if err != nil {
				retryCount = 0
			}
			if retryCount > int64(len(delays)) {
				return
			}
			delay := delays[retryCount]
			retryCount++

			fmt.Println("msg header", m.Header)
			m.Header.Add("X-Retry-Count", strconv.FormatInt(retryCount, 10))
			_ = m.NakWithDelay(delay)
		}
	})

	if err != nil {
		qm.log.Errorf("nc.QueueSubscribe: %s", err.Error())
	}
}

func (qm *QueueManager) GetRemote(subj string) *Queue {
	queue := qm.getQueue(subj)
	if queue == nil {
		queue = qm.initQueue(subj)
	}

	return queue
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
