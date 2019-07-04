package mono

import (
	"context"

	"github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/scheduler"
)

type monoScheduleOn struct {
	source Mono
	sc     scheduler.Scheduler
}

func (m monoScheduleOn) DoOnNext(fn rs.FnOnNext) Mono {
	return newMonoPeek(m, peekNext(fn))
}

func (m monoScheduleOn) Block(ctx context.Context) (interface{}, error) {
	return toBlock(ctx, m)
}

func (m monoScheduleOn) FlatMap(f flatMapper) Mono {
	return newMonoFlatMap(m, f)
}

func (m monoScheduleOn) SubscribeOn(sc scheduler.Scheduler) Mono {
	return newMonoScheduleOn(m, sc)
}

func (m monoScheduleOn) Subscribe(ctx context.Context, options ...rs.SubscriberOption) {
	m.SubscribeRaw(ctx, rs.NewSubscriber(options...))
}

func (m monoScheduleOn) SubscribeRaw(ctx context.Context, s rs.Subscriber) {
	w := m.sc.Worker()
	w.Do(func() {
		defer func() {
			_ = w.Close()
		}()
		m.source.SubscribeRaw(ctx, s)
	})
}

func (m monoScheduleOn) Filter(f rs.Predicate) Mono {
	return newMonoFilter(m, f)
}

func (m monoScheduleOn) Map(t rs.Transformer) Mono {
	return newMonoMap(m, t)
}

func newMonoScheduleOn(s Mono, sc scheduler.Scheduler) Mono {
	return monoScheduleOn{
		source: s,
		sc:     sc,
	}
}