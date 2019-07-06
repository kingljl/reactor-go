package mono

import (
	"context"
	"time"

	"github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/scheduler"
)

type wrapper struct {
	rs.RawPublisher
}

func (p wrapper) Subscribe(ctx context.Context, options ...rs.SubscriberOption) {
	p.SubscribeWith(ctx, rs.NewSubscriber(options...))
}

func (p wrapper) SwitchIfEmpty(alternative Mono) Mono {
	return wrap(newMonoSwitchIfEmpty(p, alternative))
}

func (p wrapper) Filter(f rs.Predicate) Mono {
	return wrap(newMonoFilter(p, f))
}

func (p wrapper) Map(t rs.Transformer) Mono {
	return wrap(newMonoMap(p, t))
}

func (p wrapper) FlatMap(mapper flatMapper) Mono {
	return wrap(newMonoFlatMap(p, mapper))
}

func (p wrapper) SubscribeOn(sc scheduler.Scheduler) Mono {
	return wrap(newMonoScheduleOn(p, sc))
}

func (p wrapper) DoOnNext(fn rs.FnOnNext) Mono {
	return wrap(newMonoPeek(p, peekNext(fn)))
}

func (p wrapper) DoOnError(fn rs.FnOnError) Mono {
	return wrap(newMonoPeek(p, peekError(fn)))
}

func (p wrapper) DoOnComplete(fn rs.FnOnComplete) Mono {
	return wrap(newMonoPeek(p, peekComplete(fn)))
}

func (p wrapper) DoOnCancel(fn rs.FnOnCancel) Mono {
	return wrap(newMonoPeek(p, peekCancel(fn)))
}

func (p wrapper) DoFinally(fn rs.FnOnFinally) Mono {
	return wrap(newMonoDoFinally(p, fn))
}

func (p wrapper) DelayElement(delay time.Duration) Mono {
	return wrap(newMonoDelayElement(p, delay, scheduler.Elastic()))
}

func (p wrapper) Block(ctx context.Context) (v interface{}, err error) {
	ch := make(chan struct {
		e error
		v interface{}
	}, 1)
	p.DoFinally(func(signal rs.Signal) {
		close(ch)
	}).Subscribe(ctx,
		rs.OnNext(func(s rs.Subscription, v interface{}) {
			ch <- struct {
				e error
				v interface{}
			}{e: nil, v: v}
		}),
		rs.OnError(func(e error) {
			ch <- struct {
				e error
				v interface{}
			}{e: e}
		}),
	)
	vv, ok := <-ch
	if !ok {
		return
	}
	v, err = vv.v, vv.e
	return
}

func wrap(r rs.RawPublisher) Mono {
	return wrapper{r}
}