package mono

import (
	"context"

	rs "github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/scheduler"
)

type filterSubscriber struct {
	s rs.Subscriber
	f rs.Predicate
}

func (f filterSubscriber) OnComplete() {
	f.s.OnComplete()
}

func (f filterSubscriber) OnError(err error) {
	f.s.OnError(err)
}

func (f filterSubscriber) OnNext(s rs.Subscription, v interface{}) {
	if f.f(v) {
		f.s.OnNext(s, v)
	}
}

func (f filterSubscriber) OnSubscribe(s rs.Subscription) {
	f.s.OnSubscribe(s)
}

type monoFilter struct {
	s Mono
	f rs.Predicate
}

func (m monoFilter) DoOnNext(fn rs.FnOnNext) Mono {
	return newMonoPeek(m, peekNext(fn))
}

func (m monoFilter) Block(ctx context.Context) (interface{}, error) {
	return toBlock(ctx, m)
}

func (m monoFilter) FlatMap(f flatMapper) Mono {
	return newMonoFlatMap(m, f)
}

func (m monoFilter) SubscribeOn(sc scheduler.Scheduler) Mono {
	return newMonoScheduleOn(m, sc)
}

func (m monoFilter) Filter(p rs.Predicate) Mono {
	return newMonoFilter(m, p)
}

func (m monoFilter) Subscribe(ctx context.Context, options ...rs.SubscriberOption) {
	m.SubscribeRaw(ctx, rs.NewSubscriber(options...))
}

func (m monoFilter) SubscribeRaw(ctx context.Context, s rs.Subscriber) {
	m.s.SubscribeRaw(ctx, filterSubscriber{
		s: s,
		f: m.f,
	})
}

func (m monoFilter) Map(t rs.Transformer) Mono {
	return newMonoMap(m, t)
}

func newMonoFilter(s Mono, f rs.Predicate) Mono {
	return monoFilter{
		s: s,
		f: f,
	}
}