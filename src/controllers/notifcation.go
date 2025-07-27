package controllers

import (
	"pg_bridge_go/global_var"
	"pg_bridge_go/helper"
	"pg_bridge_go/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandlePostNotificationFromPG(c *fiber.Ctx) error {
	VendorCode := c.Params("vendorcode")
	VendorPrefix := strings.SplitN(VendorCode, "-", 2)[0]

	switch VendorPrefix {
	case global_var.PGVendor.Midtrans:
		var QueryParam MidtransNotificationStruct
		if err := c.BodyParser(&QueryParam); err != nil {
			return helper.SendResponse(fiber.StatusBadRequest, fiber.Map{"error": err.Error() + " Error BindingJSON"}, nil, c)
		}

		if QueryParam.TransactionStatus == "settlement" || QueryParam.TransactionStatus == "capture" {
			err := models.UpdatePGTransactionStatus(QueryParam.OrderID, QueryParam.TransactionStatus, QueryParam.PaymentType, "midtrans-callback", global_var.DB)
			if err != nil {
				return helper.SendResponse(fiber.StatusInternalServerError, fiber.Map{"error": "Failed to update status: " + err.Error()}, nil, c)
			}
		}
	}

	return helper.SendResponse(fiber.StatusOK, fiber.Map{"message": "Notification handled"}, nil, c)
}
