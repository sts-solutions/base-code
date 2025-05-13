package ccnatsconsumer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sts-solutions/base-code/ccerrors"
	"github.com/sts-solutions/base-code/ccmetrics"
	"github.com/sts-solutions/base-code/ccmiddlewares/cccorrelation"
	"github.com/sts-solutions/base-code/ccotel/ccotelnats"
	"github.com/sts-solutions/base-code/ccworkerpool"

	"github.com/nats-io/nats.go"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
)

type JetStreamContext struct {
	nats.JetStreamContext
}

type NatsConsumer struct {
	SubDefinition    []SubscriptionDefinition
	natsMetricts     ccmetrics.MetricsNats
	js               nats.JetStreamContext
	wg               *sync.WaitGroup
	stopSubcriptions func()
}

type SubscriptionDefinition struct {
	Subject       string
	Stream        string
	ConsumerName  string
	MaxConcurrent int
	Handler       func(ctx context.Context, subj string, header nats.Header, data []byte) error
}

func New(js nats.JetStreamContext, natsMetricts ccmetrics.MetricsNats) *NatsConsumer {
	return &NatsConsumer{
		js:            js,
		wg:            &sync.WaitGroup{},
		SubDefinition: make([]SubscriptionDefinition, 0),
		natsMetricts:  natsMetricts,
	}
}

func (c *NatsConsumer) Start(ctx context.Context, shutdownCallBack func()) {
	ctx, c.stopSubcriptions = context.WithCancel(ctx)
	for _, subDef := range c.SubDefinition {
		c.wg.Add(1)
		go c.subcribe(ctx, subDef, shutdownCallBack)
	}
}

func (c *NatsConsumer) Stop(ctx context.Context) {
	log.Info(ctx, "stopping nats consumer")
	c.stopSubcriptions()
	c.wg.Wait()
	log.Info(ctx, "consumer stopped")
}

func (c *NatsConsumer) subcribe(ctx context.Context, subDef SubscriptionDefinition, shutdown func()) {
	defer c.wg.Done()

	subscription, err := c.js.PullSubscribe(subDef.Subject, subDef.ConsumerName, nats.Bind(subDef.Stream, subDef.ConsumerName))
	if err != nil {
		logEntry := log.Logger().WithFields(logrus.Fields{
			"stream":   subDef.Stream,
			"consumer": subDef.ConsumerName,
			"subject":  subDef.Subject,
		})
		logEntry.Errorf("%v", err)
		log.Info("Nats consumer")
		shutdown()
		return
	}

	msgPool := ccworkerpool.New(250)
	defer msgPool.WaitToFinish()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fetch := msgPool.Avaliable()
			c.natsMetricts.ConcurrentProcessorInc(fetch, subDef.ConsumerName)
			if fetch < 1 {
				<-time.After(time.Millisecond * 10)
				continue
			}

			msgs, err := subscription.Fetch(1)
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					continue
				}
				log.Logger().WithFields(logrus.Fields{
					"stream":   subDef.Stream,
					"consumer": subDef.ConsumerName,
					"subject":  subDef.Subject,
				}).Errorf("%v", err)

				log.Logger().Info("error feching messages")
				shutdown()
				return
			}

			for _, m := range msgs {
				msgPool.Add()

				meta, err := m.Metadata()
				if err == nil {
					c.natsMetricts.MessageLag(meta.Timestamp, []string{subDef.ConsumerName})
				}

				ctx, span := ccotelnats.StartSubscriberSpan(context.Background(), m, subDef.ConsumerName)

				ctx = context.WithValue(ctx, cccorrelation.Key,
					m.Header.Get(cccorrelation.Key.String()))

				go func(ctx context.Context, msg *nats.Msg) {
					defer msgPool.Remove()
					defer span.End()

					var handlerError error
					defer func() {
						if err := recover(); err != nil {
							handlerError = ccerrors.NewDebugTrackError("panic caught", nil, nil)
						}

						if handlerError != nil {
							if errors.Is(handlerError, &ccerrors.TransientError{}) {
								log.Logger().WithFields(logrus.Fields{
									"msg_headers": msg.Header,
									"msg_subject": msg.Subject,
									"msg_data":    string(msg.Data),
								}).Errorf("%v", handlerError)
								log.Logger().Info("nats err executin a transaction")

								if err := msg.Nak(); err != nil {
									log.Logger().
										Errorf("%v", err)
									log.Logger().Info("nats error nak-ing msg")
								}

								return
							}

							log.Logger().WithFields(logrus.Fields{
								"msg_headers": msg.Header,
								"msg_subject": msg.Subject,
								"msg_data":    string(msg.Data),
							}).Errorf("%v", handlerError)
							log.Logger().Info("nats error handling msg")
						}

						if err := msg.Ack(); err != nil {
							log.Logger().WithFields(logrus.Fields{
								"msg_headers": msg.Header,
								"msg_subject": msg.Subject,
								"msg_data":    string(msg.Data),
							}).Errorf("%v", handlerError)
							log.Logger().Info("nats err acking message")
						}
					}()
					handlerError = subDef.Handler(ctx, subDef.Subject, msg.Header, msg.Data)
				}(ctx, m)

			}

		}
	}
}
