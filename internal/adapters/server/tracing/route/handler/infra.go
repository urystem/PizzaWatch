package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"pizza/internal/domain"
	"pizza/internal/ports"
)

type handler struct {
	use ports.TrackingUse
}

type TrackingHandle interface {
	GetOrderStatus(w http.ResponseWriter, r *http.Request)
	GetOrderHistory(w http.ResponseWriter, r *http.Request)
	GetWorkers(w http.ResponseWriter, r *http.Request)
}

func NewHandler(use ports.TrackingUse) TrackingHandle {
	return &handler{use}
}

type myErr struct {
	ErrStr string `json:"error"`
}

func errorWrite(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	msg := &myErr{
		ErrStr: err.Error(),
	}
	json.NewEncoder(w).Encode(msg)
}

func (h *handler) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	ord, err := h.use.OrderStatusUpdate(r.Context(), r.PathValue("order_number"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			errorWrite(w, http.StatusNotFound, err)
		} else {
			errorWrite(w, http.StatusInternalServerError, err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ord)
}

func (h *handler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	ord, err := h.use.GetOrderHistory(r.Context(), r.PathValue("order_number"))
	if err != nil {
		errorWrite(w, http.StatusInternalServerError, err)
		return
	}
	if len(ord) == 0 {
		errorWrite(w, http.StatusNotFound, domain.ErrNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ord)
}

func (h *handler) GetWorkers(w http.ResponseWriter, r *http.Request) {
	ord, err := h.use.GetWorkersStatus(r.Context(), 10)
	if err != nil {
		errorWrite(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ord)
}
