package middleware

import "net/http"

type Middlewarer interface {
	Handle(h http.Handler) http.Handler
}

type Middleware []Middlewarer

func New(ms ...Middlewarer) Middleware {
	return append(Middleware{}, ms...)
}

func (m Middleware) Append(ms ...Middlewarer) Middleware {
	return append(m, ms...)
}

func (m Middleware) Apply(h http.Handler) http.Handler {
	ret := h
	for i := len(m) - 1; i >= 0; i-- {
		ret = m[i].Handle(ret)
	}
	return ret
}
