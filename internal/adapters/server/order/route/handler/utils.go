package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"pizza/internal/domain"
)

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

func validateOrder(ord *domain.Order) error {
	ord.CustomerName = strings.TrimSpace(ord.CustomerName)
	if ord.CustomerName == "" {
		return fmt.Errorf("invalid name")
	}
	if len(ord.Items) == 0 {
		return fmt.Errorf("no items")
	}
	if ord.OrderType == "dinein" {
		if ord.TableNumber == nil || *ord.TableNumber == 0 {
			return fmt.Errorf("give table number")
		}
	} else if ord.OrderType == "delivery" {
		if ord.DeliveryAddr == nil {
			return fmt.Errorf("give delivery addr")
		}
		*ord.DeliveryAddr = strings.TrimSpace(*ord.DeliveryAddr)
		if len(*ord.DeliveryAddr) == 0 {
			return fmt.Errorf("invalid addr")
		}
	} else if ord.OrderType != "takeout" {
		return fmt.Errorf("invalid type order: %s", ord.OrderType)
	}
	for i, v := range ord.Items {
		v.Name = strings.TrimSpace(v.Name)
		if len(v.Name) == 0 {
			return fmt.Errorf("invalid item name, number:%d", i+1)
		}
		ord.Items[i].Name = v.Name

		if v.Price < 0 {
			return fmt.Errorf("invalid price of items, number:%d", i+1)
		}
	}
	return nil
}
