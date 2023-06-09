package util

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// https://gobyexample.com.ru/signals

type Waiter struct {
	cancel context.CancelFunc
	sigs   chan os.Signal
}

func NewWaiter(cancel context.CancelFunc) Waiter {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return Waiter{
		cancel: cancel,
		sigs:   sigs,
	}
}

func (w Waiter) Wait() {
	<-w.sigs
	w.cancel()
}

func (w Waiter) Cancel() {
	w.cancel()
}

func NewStopCtx() (context.Context, Waiter) {
	ctx, cancel := context.WithCancel(context.Background())
	return ctx, NewWaiter(cancel)
}
