package internal

import (
	"context"
	"net/http"
)

type Http struct {
	service *http.Server
}

func NewHttpService(router http.Handler) *Http {
	return &Http{
		service: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}

func (h *Http) Start() error {
	if err := h.service.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (h *Http) Stop(ctx context.Context) error {
	return h.service.Shutdown(ctx)
}
