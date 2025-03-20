package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ispindel.piwo.org/internal/handlers"
	"ispindel.piwo.org/internal/services"
	"ispindel.piwo.org/pkg/auth"
	"ispindel.piwo.org/pkg/database"
)

func main() {
	// Inicjalizacja bazy danych
	database.InitDB()

	// Inicjalizacja routera Gin
	r := gin.Default()

	// Konfiguracja szablonów HTML
	r.LoadHTMLGlob("web/templates/*")

	// Konfiguracja statycznych plików
	r.Static("/static", "./web/static")

	// Inicjalizacja serwisów
	userService := services.NewUserService()

	// Middleware do sprawdzania autentykacji
	authMiddleware := func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.Next()
			return
		}

		userID, err := auth.ValidateToken(token)
		if err != nil {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Next()
			return
		}

		user, err := userService.GetUserByID(userID)
		if err != nil {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Next()
			return
		}

		c.Set("user", user)
		c.Next()
	}

	// Inicjalizacja handlerów
	authHandler := handlers.NewAuthHandler()

	// Użyj middleware'a dla wszystkich routów
	r.Use(authMiddleware)

	// Grupa routów dla autentykacji
	auth := r.Group("/auth")
	{
		auth.GET("/login", authHandler.Login)
		auth.POST("/login", authHandler.Login)
		auth.GET("/register", authHandler.Register)
		auth.POST("/register", authHandler.Register)
		auth.GET("/logout", authHandler.Logout)
	}

	// Strona główna
	r.GET("/", func(c *gin.Context) {
		user, _ := c.Get("user")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"user": user,
		})
	})

	// Panel główny (chroniony)
	r.GET("/dashboard", func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			return
		}
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"user": user,
		})
	})

	// Uruchomienie serwera
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
} 