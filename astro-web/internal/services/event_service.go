package services

import (
	"context"
	"errors"
	"strings"

	"astro-web/internal/repository"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

// CreateEvent valida las reglas de negocio e integridad antes de persistir
func (s *EventService) CreateEvent(ctx context.Context, event *repository.Event) error {
	if strings.TrimSpace(event.Title) == "" {
		return errors.New("el título del evento no puede estar vacío")
	}
	if strings.TrimSpace(event.Category) == "" {
		return errors.New("la categoría del evento no puede estar vacía")
	}
	if event.EventDate.IsZero() {
		return errors.New("la fecha del evento es obligatoria")
	}

	return s.repo.CreateEvent(event)
}

// ListEvents recupera todos los eventos ordenados por fecha
func (s *EventService) ListEvents(ctx context.Context) ([]repository.Event, error) {
	return s.repo.GetAll(ctx)
}

// GetEvent busca un evento validando que el ID sea válido
func (s *EventService) GetEvent(ctx context.Context, id int64) (*repository.Event, error) {
	if id <= 0 {
		return nil, errors.New("el id debe ser mayor a cero")
	}
	return s.repo.GetByID(id)
}

// DeleteEvent elimina un evento validando que el ID sea válido
func (s *EventService) DeleteEvent(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("el id debe ser mayor a cero")
	}
	return s.repo.Delete(id)
}
