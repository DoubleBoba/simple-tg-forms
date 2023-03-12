package function //nolint:all

import "context"

func ReadChannel[T any](ctx context.Context, channel chan T) {
	select {
	case <-ctx.Done():
		return
	case <-channel:
		return
	}
}

func AsyncFunc(function func(chan bool)) chan bool {
	done_channel := make(chan bool)
	go function(done_channel)
	return done_channel
}
