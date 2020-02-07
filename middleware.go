package middleware

import "net/http"

type Middlewarer interface {
	Handle(h http.Handler) http.Handler
}

type Middleware struct {
	middlewares []Middlewarer
}

func New(ms ...Middlewarer) *Middleware {
	return &Middleware{middlewares: ms}
}

func (m *Middleware) Append(ms ...Middlewarer) *Middleware {
	return &Middleware{
		middlewares: append(m.middlewares, ms...),
	}
}

func (m *Middleware) Apply(h http.Handler) http.Handler {
	ret := h
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		ret = m.middlewares[i].Handle(ret)
	}
	return ret
}
