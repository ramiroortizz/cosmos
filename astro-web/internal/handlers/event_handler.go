package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"astro-web/internal/repository"
	"astro-web/internal/services"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(service *services.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// RegisterRoutes registra las rutas HTTP para eventos usando el mux nativo de Go
func (h *EventHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/events", h.ListEvents)
	mux.HandleFunc("POST /api/events", h.CreateEvent)
	mux.HandleFunc("GET /api/events/{id}", h.GetEvent)
	mux.HandleFunc("DELETE /api/events/{id}", h.DeleteEvent)
}

func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.ListEvents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event repository.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateEvent(r.Context(), &event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	event, err := h.service.GetEvent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteEvent(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
