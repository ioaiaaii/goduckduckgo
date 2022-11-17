// Package config provides functions that allow to run the service. It provides default values but also overwrite them from the env vars.
package http

import (
	"time"
)

type options struct {
	HTTPPort                  string
	Env                       string
	HTTPServerTimeout         time.Duration
	HTTPServerShutdownTimeout time.Duration
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func HTTPPort(t string) Option {
	return optionFunc(func(o *options) {
		o.HTTPPort = t
	})
}

func Env(s string) Option {
	return optionFunc(func(o *options) {
		o.Env = s
	})
}

func HTTPServerTimeout(mux time.Duration) Option {
	return optionFunc(func(o *options) {
		o.HTTPServerTimeout = mux
	})
}

func HTTPServerShutdownTimeout(mux time.Duration) Option {
	return optionFunc(func(o *options) {
		o.HTTPServerTimeout = mux
	})
}
