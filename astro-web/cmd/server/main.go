package main

import (
	"database/sql"
	"log"
	"net/http"

	"astro-web/internal/handlers"
	"astro-web/internal/repository"
	"astro-web/internal/services"

	_ "modernc.org/sqlite"
)

func main() {
	// 1. Abrir la conexión a la base de datos SQLite
	db, err := sql.Open("sqlite", "./astro.db")
	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}
	defer db.Close()

	// 2. Verificar conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	// 3. Instanciar Capas (Repository -> Service -> Handler)
	eventRepo := repository.NewEventRepository(db)
	if err := eventRepo.InitTable(); err != nil {
		log.Fatalf("Error al inicializar la tabla: %v", err)
	}

	eventService := services.NewEventService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventService)

	// 4. Configurar Router HTTP nativo
	mux := http.NewServeMux()
	eventHandler.RegisterRoutes(mux)

	// 5. Levantar servidor HTTP
	port := ":8080"
	log.Printf("🚀 Servidor corriendo en http://localhost%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Error en el servidor HTTP: %v", err)
	}
}
