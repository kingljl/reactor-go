package mono

import (
	"context"
	"errors"
	"time"

	"github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/scheduler"
)

var errNotProcessor = errors.New("publisher is not a Processor")

type wrapper struct {
	reactor.RawPublisher
}

func (p wrapper) Subscribe(ctx context.Context, options ...reactor.SubscriberOption) {
	p.SubscribeWith(ctx, reactor.NewSubscriber(options...))
}

func (p wrapper) SwitchIfEmpty(alternative Mono) Mono {
	return wrap(newMonoSwitchIfEmpty(p.RawPublisher, alternative))
}

func (p wrapper) Filter(f reactor.Predicate) Mono {
	return wrap(newMonoFilter(p.RawPublisher, f))
}

func (p wrapper) Map(t reactor.Transformer) Mono {
	return wrap(newMonoMap(p.RawPublisher, t))
}

func (p wrapper) FlatMap(mapper flatMapper) Mono {
	return wrap(newMonoFlatMap(p.RawPublisher, mapper))
}

func (p wrapper) SubscribeOn(sc scheduler.Scheduler) Mono {
	return wrap(newMonoScheduleOn(p.RawPublisher, sc))
}

func (p wrapper) DoOnNext(fn reactor.FnOnNext) Mono {
	return wrap(newMonoPeek(p.RawPublisher, peekNext(fn)))
}

func (p wrapper) DoOnError(fn reactor.FnOnError) Mono {
	return wrap(newMonoPeek(p.RawPublisher, peekError(fn)))
}

func (p wrapper) DoOnComplete(fn reactor.FnOnComplete) Mono {
	return wrap(newMonoPeek(p.RawPublisher, peekComplete(fn)))
}

func (p wrapper) DoOnCancel(fn reactor.FnOnCancel) Mono {
	return wrap(newMonoPeek(p.RawPublisher, peekCancel(fn)))
}

func (p wrapper) DoOnDiscard(fn reactor.FnOnDiscard) Mono {
	return wrap(newMonoContext(p.RawPublisher, withContextDiscard(fn)))
}

func (p wrapper) DoFinally(fn reactor.FnOnFinally) Mono {
	return wrap(newMonoDoFinally(p.RawPublisher, fn))
}

func (p wrapper) DoOnSubscribe(fn reactor.FnOnSubscribe) Mono {
	return wrap(newMonoPeek(p.RawPublisher, peekSubscribe(fn)))
}

func (p wrapper) DelayElement(delay time.Duration) Mono {
	return wrap(newMonoDelayElement(p.RawPublisher, delay, scheduler.Elastic()))
}

func (p wrapper) Block(ctx context.Context) (Any, error) {
	ch := make(chan Any, 1)
	p.
		DoFinally(func(signal reactor.SignalType) {
			close(ch)
		}).
		Subscribe(ctx,
			reactor.OnNext(func(v Any) error {
				ch <- v
				return nil
			}),
			reactor.OnError(func(e error) {
				ch <- e
			}),
		)
	v, ok := <-ch
	if !ok {
		return nil, nil
	}
	if err, ok := v.(error); ok {
		return nil, err
	}
	return v, nil
}

func (p wrapper) Success(v Any) {
	p.mustProcessor().Success(v)
}

func (p wrapper) Error(e error) {
	p.mustProcessor().Error(e)
}

func (p wrapper) mustProcessor() *processor {
	pp, ok := p.RawPublisher.(*processor)
	if !ok {
		panic(errNotProcessor)
	}
	return pp
}

func wrap(r reactor.RawPublisher) wrapper {
	return wrapper{r}
}
