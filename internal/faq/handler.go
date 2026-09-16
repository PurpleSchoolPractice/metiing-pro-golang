package faq

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FAQHandler struct {
	service *FAQService
}

func NewFAQHandler(mux *chi.Mux, service *FAQService) *FAQHandler {
	handler := &FAQHandler{service: service}
	mux.Handle("GET /faq", handler.GetAll())

	return handler
}
func (h *FAQHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		faqs, err := h.service.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(faqs)
	}
}
