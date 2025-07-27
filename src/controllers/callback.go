package controllers

import (
	"fmt"
	"pg_bridge_go/db_var"
	"pg_bridge_go/global_var"
	"pg_bridge_go/helper"
	"pg_bridge_go/logger"
	"pg_bridge_go/models"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func PaymentCallback(c *fiber.Ctx) error {
	orderID := c.Query("order_id")
	// status := c.Query("transaction_status")
	// VendorCode := c.Params("vendorcode")

	LoadStatus := true

	var TransactionData db_var.PaymentGatewayTransactionT
	if err := global_var.DB.Where("order_id = ?", orderID).First(&TransactionData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.SendResponse(fiber.StatusBadRequest, "Credential not found", nil, c)
		}
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	var credential db_var.PaymentGatewayCredentialT
	if err := global_var.DB.Where("code = ? AND user_code = ?", TransactionData.Vendor, TransactionData.UserCode).First(&credential).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.SendResponse(fiber.StatusBadRequest, "Credential not found", nil, c)
		}
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	err := global_var.DB.Transaction(func(tx *gorm.DB) error {
		var err error

		TransactionStatus, PaymentType, err := SendGetPaymentStatusToMidtrans(orderID, credential)
		if err != nil {
			return err
		}

		if TransactionStatus == "capture" || TransactionStatus == "settlement" {
			err = models.UpdatePGTransactionStatus(orderID, TransactionStatus, PaymentType, "midtrans-callback", tx)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("transaction not success")
		}
		return nil
	})

	if err != nil {
		LoadStatus = false
		logger.Error("Internal server error",
			zap.Any("status_code", fiber.StatusInternalServerError),
			zap.Any("message", err),
		)
	}

	data := fiber.Map{"order_id": orderID}
	if LoadStatus {
		return c.Render("callback_success.html", data)
	} else {
		return c.Render("callback_failed.html", data)
	}
}
