package controllers

import (
	"errors"
	"net/http"
	"pg_bridge_go/db_var"
	"pg_bridge_go/global_var"
	"pg_bridge_go/helper"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterHandler(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.SendResponse(http.StatusBadRequest, "Invalid data format", nil, c)
	}

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
