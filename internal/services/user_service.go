package services

import (
	"errors"
	"time"

	"ispindel.piwo.org/internal/models"
	"ispindel.piwo.org/pkg/auth"
	"ispindel.piwo.org/pkg/database"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(name, email, password, ip string) error {
	// Sprawdź czy użytkownik już istnieje
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.New("użytkownik o podanym adresie email już istnieje")
	}

	// Hashuj hasło
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	// Utwórz nowego użytkownika
	user := models.User{
		Name:           name,
		Email:          email,
		Password:       hashedPassword,
		RegistrationIP: ip,
		IsActive:       true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(email, password, ip string) (string, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("nieprawidłowy email lub hasło")
	}

	// Sprawdź czy konto jest zablokowane
	if !user.IsActive {
		return "", errors.New("konto jest zablokowane")
	}

	if user.LockedUntil.After(time.Now()) {
		return "", errors.New("konto jest tymczasowo zablokowane")
	}

	// Sprawdź hasło
	if !auth.CheckPassword(password, user.Password) {
		// Zwiększ licznik nieudanych prób
		user.FailedLogins++
		if user.FailedLogins >= 5 {
			user.LockedUntil = time.Now().Add(15 * time.Minute)
		}
		database.DB.Save(&user)
		return "", errors.New("nieprawidłowy email lub hasło")
	}

	// Resetuj licznik nieudanych prób
	user.FailedLogins = 0
	user.LastLoginAt = time.Now()
	database.DB.Save(&user)

	// Generuj token JWT
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
} 