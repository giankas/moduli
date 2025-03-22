package main

import (
	"net/http"
	"path/filepath"

	"github.com/giankas/moduli/backend/internal/auth"
	localmw "github.com/giankas/moduli/backend/internal/middleware"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Configura il middleware per le sessioni con un CookieStore (usa una chiave sicura in produzione)
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret-key"))))
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())

	// Rotte pubbliche
	e.GET("/", func(c echo.Context) error {
		// Renderizza una pagina iniziale pubblica (puoi personalizzarla in views/index.html)
		return c.File(filepath.Join("views", "index.html"))
	})
	e.POST("/register", auth.RegisterHandler)
	e.POST("/login", auth.LoginHandler)

	// Rotte protette per la dashboard e videolezioni
	dashboardGroup := e.Group("/dashboard")
	dashboardGroup.Use(localmw.SessionAuthMiddleware)
	// La dashboard renderizza la pagina con il menu laterale e la sezione "Videolezione"
	dashboardGroup.GET("", func(c echo.Context) error {
		return c.File(filepath.Join("views", "dashboard.html"))
	})
	// Endpoint per programmare una videolezione (accessibile solo ai docenti)
	dashboardGroup.POST("/videolezioni", auth.ScheduleVideoLessonHandler)
	// Endpoint per recuperare i dettagli di una videolezione
	dashboardGroup.GET("/videolezioni/:id", auth.GetVideoLessonHandler)

	// Endpoint WebSocket per signaling WebRTC (placeholder)
	e.GET("/ws/signaling", signalingHandler)

	// File statici (CSS, JS, immagini)
	e.Static("/static", "static")

	e.Logger.Fatal(e.Start(":8080"))
}

// signalingHandler è un placeholder per il signaling WebRTC (da implementare)
func signalingHandler(c echo.Context) error {
	// Qui inserisci la logica per il signaling via WebSocket (placeholder)
	return c.String(http.StatusOK, "Signaling endpoint (da implementare)")
}
