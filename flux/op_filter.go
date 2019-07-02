package flux

import (
	"context"

	rs "github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/scheduler"
)

type filterSubscriber struct {
	source    rs.Subscriber
	predicate rs.Predicate
}

func (s filterSubscriber) OnComplete() {
	s.source.OnComplete()
}

func (s filterSubscriber) OnError(err error) {
	s.source.OnError(err)
}

func (s filterSubscriber) OnNext(su rs.Subscription, in interface{}) {
	if s.predicate(in) {
		s.source.OnNext(su, in)
	}
}

func (s filterSubscriber) OnSubscribe(ss rs.Subscription) {
	s.source.OnSubscribe(ss)
}

func (filterSubscriber) Raw() rs.RawSubscriber {
	panic("implement me")
}

func newFilterSubscriber(s rs.Subscriber, p rs.Predicate) rs.Subscriber {
	return filterSubscriber{
		source:    s,
		predicate: p,
	}
}

type fluxFilter struct {
	source    Flux
	predicate rs.Predicate
}

func (f fluxFilter) Subscribe(ctx context.Context, s rs.Subscriber) rs.Disposable {
	f.source.Subscribe(ctx, newFilterSubscriber(s, f.predicate))
	return nil
}

func (f fluxFilter) Filter(p rs.Predicate) Flux {
	return newFluxFilter(f, p)
}

func (f fluxFilter) Map(t rs.Transformer) Flux {
	return newFluxMap(f, t)
}

func (f fluxFilter) SubscribeOn(sc scheduler.Scheduler) Flux {
	return newFluxSubscribeOn(f, sc)
}

func newFluxFilter(source Flux, predicate rs.Predicate) Flux {
	return fluxFilter{
		source:    source,
		predicate: predicate,
	}
}