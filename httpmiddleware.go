package httpmiddleware

import "github.com/julienschmidt/httprouter"

type Middleware func(httprouter.Handle) httprouter.Handle
type Middlewares []Middleware

type IHandler interface {
	Count() int
	Use(m ...Middleware) IHandler
	Wrap(handler httprouter.Handle) httprouter.Handle
}

type Handler struct {
	middlewares Middlewares
}

func New() IHandler {
	return &Handler{
		middlewares: Middlewares{},
	}
}

func (h Handler) Count() int {
	return len(h.middlewares)
}

func (h *Handler) Use(m ...Middleware) IHandler {
	h.middlewares = append(m, h.middlewares...)
	return h
}

func (h *Handler) Wrap(handler httprouter.Handle) httprouter.Handle {
	l := len(h.middlewares)
	if l == 0 {
		return handler
	}

	var result httprouter.Handle
	result = h.middlewares[l-1](handler) // Wrap the initial middleware

	// Backward traverse to wrap the middleware handler one by one...
	for i := l - 2; i >= 0; i-- {
		result = h.middlewares[i](result)
	}

	return result
}
