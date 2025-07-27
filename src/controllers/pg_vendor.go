package controllers

import (
	"fmt"
	"pg_bridge_go/config"
	"pg_bridge_go/db_var"
	"pg_bridge_go/global_var"
	"pg_bridge_go/helper"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePaymentGatewayCredential(c *fiber.Ctx) error {
	type Request struct {
		Vendor           string `json:"vendor" binding:"required"`
		GatewayName      string `json:"gateway_name" binding:"required"`
		APIKey           string `json:"api_key" binding:"required"`
		APISecret        string `json:"api_secret"`
		MerchantID       string `json:"merchant_id"`
		CallbackURL      string `json:"callback_url"`
		CallbackRedirect int    `json:"callback_redirect"`
		Mode             string `json:"mode"`
	}

	var input Request
	if err := c.BodyParser(&input); err != nil {
		return helper.SendResponse(fiber.StatusBadRequest, nil, nil, c)
	}

	// Validate vendor
	vendorKey := strings.ToLower(input.Vendor)
	prefix := ""
	switch vendorKey {
	case "midtrans":
		prefix = global_var.PGVendor.Midtrans
	case "xendit":
		prefix = global_var.PGVendor.Xendit
	case "hitpay":
		prefix = global_var.PGVendor.HitPay
	case "doku":
		prefix = global_var.PGVendor.Doku
	default:
		return helper.SendResponse(fiber.StatusBadRequest, "Invalid vendor", nil, c)
	}

	ApiKeyEn, err := helper.Encrypt(input.APIKey, config.MasterKey)
	if err != nil {
		return helper.SendResponse(fiber.StatusInternalServerError, fmt.Sprintf("%v", err), nil, c)
	}
	ApiSecretEn, err := helper.Encrypt(input.APISecret, config.MasterKey)
	if err != nil {
		return helper.SendResponse(fiber.StatusInternalServerError, fmt.Sprintf("%v", err), nil, c)
	}
	MerchantIDEn, err := helper.Encrypt(input.MerchantID, config.MasterKey)
	if err != nil {
		return helper.SendResponse(fiber.StatusInternalServerError, fmt.Sprintf("%v", err), nil, c)
	}

	credential := db_var.PaymentGatewayCredentialT{
		GatewayName:      input.GatewayName,
		APIKey:           ApiKeyEn,
		APISecret:        ApiSecretEn,
		MerchantID:       MerchantIDEn,
		CallbackURL:      input.CallbackURL,
		CallbackRedirect: input.CallbackRedirect,
		Mode:             input.Mode,
		UserCode:         helper.GetUsernameFiber(c),
		CreatedBy:        helper.GetUsernameFiber(c),
	}

	if err := global_var.DB.Create(&credential).Error; err != nil {
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	// Generate Code with vendor prefix
	code := fmt.Sprintf("%s-%d", prefix, credential.ID)

	if err := global_var.DB.Model(&credential).Update("code", code).Error; err != nil {
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	credential.Code = code

	return helper.SendResponse(fiber.StatusOK, "", credential, c)
}

func GetPaymentGatewayCredential(c *fiber.Ctx) error {
	code := c.Params("code")

	var credential db_var.PaymentGatewayCredentialT
	if err := global_var.DB.Where("code = ? AND user_code = ?", code, helper.GetUsernameFiber(c)).First(&credential).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.SendResponse(fiber.StatusBadRequest, "Credential not found", nil, c)
		}
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	credential.APIKey, _ = helper.Decrypt(credential.APIKey, config.MasterKey)
	credential.APISecret, _ = helper.Decrypt(credential.APISecret, config.MasterKey)
	credential.MerchantID, _ = helper.Decrypt(credential.MerchantID, config.MasterKey)

	return helper.SendResponse(fiber.StatusOK, "", credential, c)
}

func GetAllPaymentGatewayCredential(c *fiber.Ctx) error {
	var credential []db_var.PaymentGatewayCredentialT
	if err := global_var.DB.Where("user_code = ?", helper.GetUsernameFiber(c)).Find(&credential).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.SendResponse(fiber.StatusBadRequest, "Credential not found", nil, c)
		}
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	for i := range credential {
		credential[i].APIKey, _ = helper.Decrypt(credential[i].APIKey, config.MasterKey)
		credential[i].APISecret, _ = helper.Decrypt(credential[i].APISecret, config.MasterKey)
		credential[i].MerchantID, _ = helper.Decrypt(credential[i].MerchantID, config.MasterKey)
	}

	return helper.SendResponse(fiber.StatusOK, "", credential, c)
}

func UpdatePaymentGatewayCredential(c *fiber.Ctx) error {
	type Request struct {
		GatewayName      string `json:"gateway_name" binding:"required"`
		APIKey           string `json:"api_key" binding:"required"`
		APISecret        string `json:"api_secret"`
		MerchantID       string `json:"merchant_id"`
		CallbackURL      string `json:"callback_url"`
		CallbackRedirect int    `json:"callback_redirect"`
		Mode             string `json:"mode"`
	}

	code := c.Params("code")

	var credential db_var.PaymentGatewayCredentialT
	if err := global_var.DB.Where("code = ? AND user_code = ?", code, helper.GetUsernameFiber(c)).First(&credential).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.SendResponse(fiber.StatusBadRequest, "Credential not found", nil, c)
		}
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	var input Request
	if err := c.BodyParser(&input); err != nil {
		return helper.SendResponse(fiber.StatusBadRequest, nil, nil, c)
	}

	credential.GatewayName = input.GatewayName
	credential.APIKey, _ = helper.Encrypt(input.APIKey, config.MasterKey)
	credential.APISecret, _ = helper.Encrypt(input.APISecret, config.MasterKey)
	credential.MerchantID, _ = helper.Encrypt(input.MerchantID, config.MasterKey)
	credential.CallbackURL = input.CallbackURL
	credential.CallbackRedirect = input.CallbackRedirect
	credential.Mode = input.Mode
	credential.UpdatedAt = time.Now()
	credential.UpdatedBy = helper.GetUsernameFiber(c)

	if err := global_var.DB.Save(&credential).Error; err != nil {
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	return helper.SendResponse(fiber.StatusOK, "", credential, c)
}

func DeletePaymentGatewayCredential(c *fiber.Ctx) error {
	code := c.Params("code")

	if err := global_var.DB.Where("code = ? AND user_code = ?", code, helper.GetUsernameFiber(c)).Delete(&db_var.PaymentGatewayCredentialT{}).Error; err != nil {
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	return helper.SendResponse(fiber.StatusOK, "Credential deleted", nil, c)
}
