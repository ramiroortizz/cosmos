package repository

import (
	"context"
	"database/sql"
	"time"
)

type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	Category    string    `json:"category"`
}

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}
func (r *EventRepository) InitTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS events (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT,
        event_date DATETIME NOT NULL,
        category TEXT NOT NULL
    );
    `
	_, err := r.db.Exec(query)
	return err
}
func (r *EventRepository) CreateEvent(event *Event) error {
	query := `INSERT INTO events (title, description, event_date, category) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, event.Title, event.Description, event.EventDate, event.Category)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	event.ID = id
	return nil
}

// recuperar todos los eventos
func (r *EventRepository) GetAll(ctx context.Context) ([]Event, error) {
	query := `SELECT id, title, description, event_date, category FROM events ORDER BY event_date ASC;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.EventDate, &e.Category); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

// Obtener un evento por su ID
func (r *EventRepository) GetByID(id int64) (*Event, error) {
	query := `SELECT id, title, description, event_date, category FROM events WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var e Event
	if err := row.Scan(&e.ID, &e.Title, &e.Description, &e.EventDate, &e.Category); err != nil {
		return nil, err
	}
	return &e, nil
}

// Eliminar un evento por su ID
func (r *EventRepository) Delete(id int64) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
