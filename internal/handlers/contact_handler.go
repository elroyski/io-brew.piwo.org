package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"ispindel.piwo.org/internal/models"
	"ispindel.piwo.org/pkg/mailer"
)

// ContactHandler obsługuje żądania związane z formularzem kontaktowym
type ContactHandler struct{}

func NewContactHandler() *ContactHandler {
	return &ContactHandler{}
}

// SendMessage obsługuje wysyłanie wiadomości kontaktowej
func (h *ContactHandler) SendMessage(c *gin.Context) {
	// Pobierz zalogowanego użytkownika
	user, exists := c.Get("user")
	if !exists {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"error": "Nie jesteś zalogowany",
		})
		return
	}
	userModel := user.(*models.User)

	// Pobierz dane z formularza
	subject := c.PostForm("subject")
	message := c.PostForm("message")
	userDevices := c.PostForm("userDevices")
	userFermentations := c.PostForm("userFermentations")

	// Walidacja danych
	if subject == "" || message == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Wszystkie pola są wymagane",
			"user":  userModel,
		})
		return
	}

	// Przygotuj treść wiadomości
	fullMessage := fmt.Sprintf(`
Od: %s (%s)
Temat: %s

Urządzenia użytkownika: %s
Fermentacje użytkownika: %s

Treść:
%s
`, userModel.Name, userModel.Email, subject, userDevices, userFermentations, message)

	// Wyślij wiadomość do administratora
	err := mailer.SendEmail("elroyski@gmail.com", "Zgłoszenie z iSpindel: "+subject, fullMessage)
	if err != nil {
		log.Printf("Błąd wysyłania wiadomości kontaktowej: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Nie udało się wysłać wiadomości. Spróbuj ponownie później.",
			"user":  userModel,
		})
		return
	}

	// Przekieruj z powrotem na dashboard z komunikatem o sukcesie
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"user":              userModel,
		"contactSuccess":    true,
		"contactSuccessMsg": "Twoja wiadomość została wysłana. Dziękujemy za kontakt!",
		"userDevices":       userDevices,
		"userFermentations": userFermentations,
		"isAdmin":           userModel.Email == "elroyski@gmail.com",
	})
}
