package handler

import (
	"encoding/json"
	"net/http"

	"pizza/internal/domain"
	"pizza/internal/ports"
)

type handler struct {
	use ports.OrderUseCase
}

type OrderHandle interface {
	CreateOrder(w http.ResponseWriter, r *http.Request)
}

func NewHandler(use ports.OrderUseCase) OrderHandle {
	return &handler{use}
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ord := new(domain.Order)
	err := json.NewDecoder(r.Body).Decode(ord)
	if err != nil {
		errorWrite(w, http.StatusBadRequest, err)
		return
	}
	err = validateOrder(ord)
	if err != nil {
		errorWrite(w, http.StatusBadRequest, err)
		return
	}
	sts, err := h.use.CreateOrder(r.Context(), ord)
	if err != nil {
		errorWrite(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sts)
}
