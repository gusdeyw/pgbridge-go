package controllers

import (
	"errors"
	"net/http"
	"pg_bridge_go/db_var"
	"pg_bridge_go/global_var"
	"pg_bridge_go/helper"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// validateInput performs input validation and sanitization
func validateInput(username, password string) error {
	// Trim whitespace
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	// Username validation
	if len(username) < 3 || len(username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}

	// Username should contain only alphanumeric characters, underscore, and dash
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, underscores, and dashes")
	}

	// Password validation
	if len(password) < 8 || len(password) > 255 {
		return errors.New("password must be between 8 and 255 characters")
	}

	// Check for at least one uppercase, one lowercase, one digit
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one digit")
	}

	return nil
}

func RegisterHandler(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.SendResponse(http.StatusBadRequest, "Invalid data format", nil, c)
	}

	// Input validation and sanitization
	if err := validateInput(req.Username, req.Password); err != nil {
		return helper.SendResponse(http.StatusBadRequest, err.Error(), nil, c)
	}

	// Sanitize username
	req.Username = strings.TrimSpace(req.Username)

	var existing db_var.UserT
	if err := global_var.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		return helper.SendResponse(http.StatusBadRequest, "Username already taken", nil, c)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.SendResponse(http.StatusInternalServerError, "Database error", nil, c)
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.SendResponse(http.StatusInternalServerError, "Password encryption failed", nil, c)
	}

	user := db_var.UserT{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := global_var.DB.Create(&user).Error; err != nil {
		return helper.SendResponse(http.StatusInternalServerError, "Database error", nil, c)
	}

	result := fiber.Map{
		"id":       user.ID,
		"username": user.Username,
	}

	return helper.SendResponse(http.StatusOK, "Registration successful", result, c)
}

func Ping(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"message": "pong",
	})
}
