package retry

import (
	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
)

const (
	defaultRetries = 1
	noRetries      = 0
)

type RetryableFunc func() (interface{}, error)

type Retryer struct {
	callback RetryableFunc
	times    int
}

func NewRetryer() *Retryer {
	return &Retryer{
		callback: nil,
		times:    defaultRetries,
	}
}

func (r *Retryer) Callback(callback RetryableFunc) *Retryer {
	r.callback = callback
	return r
}

func (r *Retryer) Times(times int) *Retryer {
	r.times = times
	return r
}

func (r *Retryer) Retry() (interface{}, error) {
	if r.times == noRetries {
		r.times = defaultRetries
	}

	if r.callback == nil {
		return nil, errutil.NewCriticalError("callback must be set to do retry")
	}

	var (
		err    error
		result interface{}
	)

	for i := 1; i <= r.times; i++ {
		result, err = r.callback()
		if err == nil {
			return result, nil
		}
	}

	return nil, err
}
