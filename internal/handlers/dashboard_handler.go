package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"ispindel.piwo.org/internal/models"
	"ispindel.piwo.org/internal/services"
)

type DashboardHandler struct {
	IspindelService     *services.IspindelService
	FermentationService *services.FermentationService
}

func NewDashboardHandler(ispindelService *services.IspindelService, fermentationService *services.FermentationService) *DashboardHandler {
	return &DashboardHandler{
		IspindelService:     ispindelService,
		FermentationService: fermentationService,
	}
}

// ServeHTTP obsługuje żądanie wyświetlenia dashboardu
func (h *DashboardHandler) Dashboard(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return
	}
	userModel := user.(*models.User)
	userID := userModel.ID

	// Pobierz urządzenia użytkownika
	userDevices, err := h.IspindelService.GetIspindelsByUserID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Nie można pobrać urządzeń użytkownika",
			"user":  userModel,
		})
		return
	}

	// Pobierz fermentacje użytkownika
	userFermentations, err := h.FermentationService.GetFermentationsByUserID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Nie można pobrać fermentacji użytkownika",
			"user":  userModel,
		})
		return
	}

	// Przygotuj dane do szablonu
	deviceStrings := make([]string, len(userDevices))
	for i, device := range userDevices {
		deviceStrings[i] = device.Name
	}

	fermentationStrings := make([]string, len(userFermentations))
	for i, fermentation := range userFermentations {
		fermentationStrings[i] = fermentation.Name
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"user":              userModel,
		"userDevices":       strings.Join(deviceStrings, ", "),
		"userFermentations": strings.Join(fermentationStrings, ", "),
		"isAdmin":           userModel.Email == "elroyski@gmail.com",
	})
}
