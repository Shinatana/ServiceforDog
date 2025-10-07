package internal

import (
	"context"
	"net/http"
)

type Http struct {
	service *http.Server
}

func NewHttpService() *Http {
	return &Http{
		service: &http.Server{
			Addr: ":8080",
		},
	}
}

func (h *Http) Start() error {
	if err := h.service.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (h *Http) AddHandlerFunc(path string, fn http.HandlerFunc) {
	http.HandleFunc(path, fn)
}

func (h *Http) Stop(ctx context.Context) error {
	return h.service.Shutdown(ctx)
}
