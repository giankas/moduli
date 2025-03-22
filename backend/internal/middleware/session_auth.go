package middleware

import (
	"net/http"

	"github.com/giankas/moduli/backend/internal/models"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// SessionAuthMiddleware controlla se la sessione contiene l'utente autenticato
func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Errore nella gestione della sessione"})
		}
		userInterface := sess.Values["user"]
		if userInterface == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Non autorizzato"})
		}
		// Opzionalmente puoi fare cast a *models.User
		_, ok := userInterface.(*models.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Sessione non valida"})
		}
		return next(c)
	}
}
