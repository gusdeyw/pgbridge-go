package controllers

import (
	"encoding/json"
	"fmt"
	"pg_bridge_go/config"
	"pg_bridge_go/db_var"
	"pg_bridge_go/global_var"
	"pg_bridge_go/helper"
	"pg_bridge_go/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PaymentRequest struct {
	OrderID         string            `json:"order_id"`
	Amount          int               `json:"amount"`
	Items           *[]PaymentItem    `json:"items,omitempty"`
	Customer        *CustomerInfo     `json:"customer,omitempty"`
	EnabledPayments *[]string         `json:"enabled_payments,omitempty"`
	Callbacks       *CallbackURLs     `json:"callbacks,omitempty"`
	Expiry          *PaymentExpiry    `json:"expiry,omitempty"`
	CustomFields    map[string]string `json:"custom_fields,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

type PaymentItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Brand    string `json:"brand,omitempty"`
	Category string `json:"category,omitempty"`
}

type CustomerInfo struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	Billing   *Address `json:"billing,omitempty"`
	Shipping  *Address `json:"shipping,omitempty"`
}

type Address struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	AddressLine string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

type CallbackURLs struct {
	Finish string `json:"finish"`
}

type PaymentExpiry struct {
	StartTime string `json:"start_time"`
	Unit      string `json:"unit"` // e.g. "minutes"
	Duration  int    `json:"duration"`
}

func HandleCreatePayment(c *fiber.Ctx) error {
	VendorCode := c.Params("vendorcode")
	var Req PaymentRequest
	var RedirectUrl string

	if err := c.BodyParser(&Req); err != nil {
		return helper.SendResponse(fiber.StatusBadRequest, nil, nil, c)
	}

	OrderID := ""
	if Req.OrderID == "" {
		OrderID = helper.GenerateOrderID(VendorCode)
	} else {
		OrderID = Req.OrderID
	}

	var credential db_var.PaymentGatewayCredentialT
	if err := global_var.DB.Where("code = ? AND user_code = ?", VendorCode, helper.GetUsernameFiber(c)).First(&credential).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.SendResponse(fiber.StatusBadRequest, "Credential not found", nil, c)
		}
		return helper.SendResponse(fiber.StatusInternalServerError, "", nil, c)
	}

	VendorPrefix := strings.SplitN(VendorCode, "-", 2)[0]

	switch VendorPrefix {
	case global_var.PGVendor.Midtrans:
		RequestBody := MidtransTransactionRequest{
			TransactionDetails: MidtransTransactionDetails{
				OrderID:     OrderID,
				GrossAmount: Req.Amount,
			},
			CreditCard: MidtransCreditCard{Secure: true},
		}

		// Only add ItemDetails if exists
		if Req.Items != nil && len(*Req.Items) > 0 {
			var itemDetails []MidtransItemDetail
			for i, v := range *Req.Items {
				item := MidtransItemDetail{
					ID:       fmt.Sprintf("%d", i+1),
					Price:    v.Price,
					Quantity: v.Quantity,
					Name:     v.Name,
				}
				if v.Brand != "" {
					item.Brand = &v.Brand
				}
				if v.Category != "" {
					item.Category = &v.Category
				}
				itemDetails = append(itemDetails, item)
			}
			RequestBody.ItemDetails = &itemDetails
		}

		// Only add CustomerDetails if exists and valid
		if Req.Customer != nil && Req.Customer.Email != "" {
			customer := MidtransCustomerDetails{
				FirstName: Req.Customer.FirstName,
				LastName:  Req.Customer.LastName,
				Email:     Req.Customer.Email,
				Phone:     Req.Customer.Phone,
			}
			if Req.Customer.Billing.Email != "" {
				customer.BillingAddress = MidtransAddress{
					FirstName:   Req.Customer.Billing.FirstName,
					LastName:    Req.Customer.Billing.LastName,
					Email:       Req.Customer.Billing.Email,
					Phone:       Req.Customer.Billing.Phone,
					Address:     Req.Customer.Billing.AddressLine,
					City:        Req.Customer.Billing.City,
					PostalCode:  Req.Customer.Billing.PostalCode,
					CountryCode: Req.Customer.Billing.CountryCode,
				}
				customer.ShippingAddress = customer.BillingAddress
			}
			RequestBody.CustomerDetails = &customer
		}

		// Only add Expiry if duration >= 20
		if Req.Expiry != nil && Req.Expiry.Unit != "" {
			RequestBody.Expiry = &MidtransExpiry{
				Unit:     Req.Expiry.Unit,
				Duration: Req.Expiry.Duration,
			}
		}

		// Add Callbacks
		var cbUrl string
		if credential.CallbackRedirect == 1 && Req.Callbacks != nil && Req.Callbacks.Finish != "" {
			cbUrl = Req.Callbacks.Finish
		} else {
			cbUrl = config.CallbackUrl + "/callback/" + VendorCode + "/payment"
		}
		RequestBody.Callbacks = &MidtransCallbacks{Finish: &cbUrl}

		// Add enabled payments
		if Req.EnabledPayments != nil {
			RequestBody.EnabledPayments = Req.EnabledPayments
		} else {
			// fallback to default list
			RequestBody.EnabledPayments = &[]string{"credit_card", "mandiri_clickpay", "cimb_clicks", "bca_klikbca", "bca_klikpay", "bri_epay", "echannel", "mandiri_ecash", "permata_va", "bca_va", "bni_va", "other_va", "gopay", "indomaret", "alfamart", "danamon_online", "akulaku"}
		}

		err := global_var.DB.Transaction(func(tx *gorm.DB) error {
			var err error

			insert := db_var.PaymentGatewayTransactionT{
				OrderID:   OrderID,
				UserCode:  helper.GetUsernameFiber(c),
				Amount:    Req.Amount,
				Vendor:    VendorCode,
				Status:    "pending",
				CreatedAt: time.Now(),
				CreatedBy: helper.GetUsernameFiber(c),
			}

			if Req.Customer != nil {
				insert.CustomerName = Req.Customer.FirstName
				insert.CustomerEmail = Req.Customer.Email
				insert.CustomerPhone = Req.Customer.Phone
			}

			err = models.InsertPGTransaction(&insert, tx)
			if err != nil {
				return err
			}

			RedirectUrl, err = SendRequestPaymentToMidtrans(RequestBody, credential)
			if err != nil {
				return err
			}

			RedirectURLByte, err := json.Marshal([]string{RedirectUrl})
			if err != nil {
				return err
			}
			insert.CallbacksJSON = datatypes.JSON(RedirectURLByte)

			err = models.UpdatePGTransactionStatus(insert.OrderID, "pending", "", helper.GetUsernameFiber(c), tx)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return helper.SendResponse(fiber.StatusInternalServerError, fmt.Sprintf("%v", err), nil, c)
		}

		qrCode, err := helper.GenerateQRCodeBase64(RedirectUrl)
		if err != nil {
			return helper.SendResponse(fiber.StatusInternalServerError, "Failed to generate QR code", nil, c)
		}

		return helper.SendResponse(fiber.StatusOK, "", fiber.Map{
			"redirect_url": RedirectUrl,
			"qr_code":      qrCode,
		}, c)
	}

	return helper.SendResponse(fiber.StatusBadRequest, "no vendor code registered yet", nil, c)
}

func HandleGetPaymentStatus(c *fiber.Ctx) error {
	VendorCode := c.Params("vendorcode")
	// Fiber handles multiple query values differently
	OrderIDs := c.Query("order_id")
	var orderIDList []string
	if OrderIDs != "" {
		orderIDList = strings.Split(OrderIDs, ",")
	}
	StartDate := c.Query("start_date")
	EndDate := c.Query("end_date")
	Username := helper.GetUsernameFiber(c)

	var transactions []db_var.PaymentGatewayTransactionT
	db := global_var.DB.Model(&db_var.PaymentGatewayTransactionT{}).
		Where("vendor = ? AND user_code = ?", VendorCode, Username)

	if len(orderIDList) > 0 {
		db = db.Where("order_id IN ?", orderIDList)
	}

	if StartDate != "" && EndDate != "" {
		db = db.Where("DATE(created_at) >= ?", StartDate)
		db = db.Where("DATE(created_at) <= ?", EndDate)
	}

	err := db.Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return helper.SendResponse(fiber.StatusInternalServerError, err.Error(), nil, c)
	}

	if len(transactions) == 0 {
		return helper.SendResponse(fiber.StatusOK, "No transaction found", nil, c)
	}

	type DataReturnStruct struct {
		OrderID        string  `json:"order_id"`
		Amount         float64 `json:"amount"`
		PaymentMethods string  `json:"payment_methods"`
		Status         string  `json:"status"`
		PaidAt         string  `json:"paid_at"`
		CreatedAt      string  `json:"created_at"`
	}

	var DataReturn []DataReturnStruct
	for _, v := range transactions {
		DataReturn = append(DataReturn, DataReturnStruct{
			OrderID:        v.OrderID,
			Amount:         float64(v.Amount),
			PaymentMethods: v.PaymentMethods,
			Status:         v.Status,
			PaidAt:         v.PaidAt.Format("2006-01-02"),
			CreatedAt:      v.CreatedAt.Format("2006-01-02"),
		})
	}

	return helper.SendResponse(fiber.StatusOK, "", DataReturn, c)
}
