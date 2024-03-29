package handler

import "github.com/light-bringer/rates-exchanger-service/internal/service"

type Handler struct {
	service *service.RatesService
}

// NewHandler returns a new Handler with the given RatesService.
func NewHandler(service *service.RatesService) *Handler {
	return &Handler{service: service}
}
