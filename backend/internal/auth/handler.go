package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/giankas/moduli/backend/internal/models"
	"github.com/labstack/echo-contrib/session"
)

var (
	// In-memory user store (da sostituire con un database)
	users      = map[string]*models.User{}
	nextUserID = 1

	// In-memory videolezioni store
	videoLessons         = map[int]*models.VideoLesson{}
	videoLessonIDCounter = 1
)

// RegisterPayload per la registrazione
type RegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // "teacher" o "student"
}

// LoginPayload per il login
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler gestisce la registrazione degli utenti
func RegisterHandler(c echo.Context) error {
	payload := new(RegisterPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Payload non valido"})
	}
	if payload.Email == "" || payload.Password == "" || (payload.Role != "teacher" && payload.Role != "student") {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Campi mancanti o ruolo non valido"})
	}
	if _, exists := users[payload.Email]; exists {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email già registrata"})
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Errore durante la registrazione"})
	}
	user := &models.User{
		ID:           nextUserID,
		Email:        payload.Email,
		PasswordHash: string(hashed),
		Role:         payload.Role,
	}
	users[payload.Email] = user
	nextUserID++
	return c.JSON(http.StatusCreated, echo.Map{"message": "Registrazione effettuata"})
}

// LoginHandler gestisce il login e salva l'utente nella sessione
func LoginHandler(c echo.Context) error {
	payload := new(LoginPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Payload non valido"})
	}
	user, exists := users[payload.Email]
	if !exists {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Credenziali errate"})
	}
	// Verifica password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Credenziali errate"})
	}

	// Salva l'utente nella sessione
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Errore nella gestione della sessione"})
	}
	sess.Values["user"] = user
	sess.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, echo.Map{"message": "Login effettuato"})
}

// ScheduleVideoLessonHandler consente al docente di programmare una videolezione
func ScheduleVideoLessonHandler(c echo.Context) error {
	// Recupera l'utente dalla sessione (verifica già effettuata dal middleware)
	sess, _ := session.Get("session", c)
	userInterface := sess.Values["user"]
	if userInterface == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Non autorizzato"})
	}
	user := userInterface.(*models.User)
	if user.Role != "teacher" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Solo i docenti possono programmare videolezioni"})
	}
	title := c.FormValue("title")
	scheduledAtStr := c.FormValue("scheduled_at")
	if title == "" || scheduledAtStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Campi mancanti"})
	}
	scheduledAt, err := time.Parse(time.RFC3339, scheduledAtStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Formato data/ora non valido"})
	}

	vl := &models.VideoLesson{
		ID:          videoLessonIDCounter,
		Title:       title,
		ScheduledAt: scheduledAt,
		TeacherID:   user.ID,
	}
	videoLessons[videoLessonIDCounter] = vl
	videoLessonIDCounter++

	return c.JSON(http.StatusCreated, vl)
}

// GetVideoLessonHandler restituisce i dettagli di una videolezione
func GetVideoLessonHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID non valido"})
	}
	vl, exists := videoLessons[id]
	if !exists {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Videolezione non trovata"})
	}
	return c.JSON(http.StatusOK, vl)
}
