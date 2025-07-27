package controllers

import (
	"encoding/json"
	"fmt"
	"pg_bridge_go/config"
	"pg_bridge_go/db_var"
	"pg_bridge_go/global_var"
	"pg_bridge_go/helper"
	"strings"
)

type MidtransNotificationStruct struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	SettlementTime    string `json:"settlement_time"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

type MidtransTransactionRequest struct {
	TransactionDetails MidtransTransactionDetails `json:"transaction_details"`
	ItemDetails        *[]MidtransItemDetail      `json:"item_details,omitempty"`
	CustomerDetails    *MidtransCustomerDetails   `json:"customer_details,omitempty"`
	EnabledPayments    *[]string                  `json:"enabled_payments,omitempty"`
	CreditCard         MidtransCreditCard         `json:"credit_card"`
	BCAVA              *MidtransBCAVA             `json:"bca_va,omitempty"`
	BNIVA              *MidtransVABank            `json:"bni_va,omitempty"`
	PermataVA          *MidtransPermataVA         `json:"permata_va,omitempty"`
	Callbacks          *MidtransCallbacks         `json:"callbacks,omitempty"`
	Expiry             *MidtransExpiry            `json:"expiry,omitempty"`
	CustomField1       *string                    `json:"custom_field1,omitempty"`
	CustomField2       *string                    `json:"custom_field2,omitempty"`
	CustomField3       *string                    `json:"custom_field3,omitempty"`
}

type MidtransTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}

type MidtransItemDetail struct {
	ID           string  `json:"id"`
	Price        int     `json:"price"`
	Quantity     int     `json:"quantity"`
	Name         string  `json:"name"`
	Brand        *string `json:"brand,omitempty"`
	Category     *string `json:"category,omitempty"`
	MerchantName *string `json:"merchant_name,omitempty"`
}

type MidtransCustomerDetails struct {
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	Email           string          `json:"email"`
	Phone           string          `json:"phone"`
	BillingAddress  MidtransAddress `json:"billing_address"`
	ShippingAddress MidtransAddress `json:"shipping_address"`
}

type MidtransAddress struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

type MidtransCreditCard struct {
	Secure        bool                 `json:"secure"`
	Bank          *string              `json:"bank,omitempty"`
	Installment   *MidtransInstallment `json:"installment,omitempty"`
	WhitelistBins *[]string            `json:"whitelist_bins,omitempty"`
}

type MidtransInstallment struct {
	Required bool             `json:"required"`
	Terms    map[string][]int `json:"terms"`
}

type MidtransBCAVA struct {
	VANumber       string               `json:"va_number"`
	SubCompanyCode string               `json:"sub_company_code"`
	FreeText       MidtransFreeTextInfo `json:"free_text"`
}

type MidtransFreeTextInfo struct {
	Inquiry []MidtransLanguageText `json:"inquiry"`
	Payment []MidtransLanguageText `json:"payment"`
}

type MidtransLanguageText struct {
	EN string `json:"en"`
	ID string `json:"id"`
}

type MidtransVABank struct {
	VANumber string `json:"va_number"`
}

type MidtransPermataVA struct {
	VANumber      string `json:"va_number"`
	RecipientName string `json:"recipient_name"`
}

type MidtransCallbacks struct {
	Finish *string `json:"finish,omitempty"`
}

type MidtransExpiry struct {
	StartTime *string `json:"start_time,omitempty"`
	Unit      string  `json:"unit"`
	Duration  int     `json:"duration"`
}

type MidtransErrorResponse struct {
	ErrorMessages []string `json:"error_messages"`
}

type MidtransSuccessResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func SendRequestPaymentToMidtrans(Data MidtransTransactionRequest, Vendor db_var.PaymentGatewayCredentialT) (string, error) {
	UrlEnvMode := global_var.PGUrlList.Midtrans.Dev
	if Vendor.Mode == "prod" {
		UrlEnvMode = global_var.PGUrlList.Midtrans.Prod
	}

	ApiKeys, err := helper.Decrypt(Vendor.APIKey, config.MasterKey)
	if err != nil {
		return "", err
	}

	Reqs := helper.RequestOptions{
		Method:      "POST",
		URL:         UrlEnvMode + "/snap/v1/transactions",
		Body:        Data,
		AuthType:    helper.AuthBasic,
		Username:    ApiKeys,
		ContentType: "application/json",
	}

	Result, HttpStatus, _, err := helper.SendRequest(Reqs)
	if err != nil {
		return "", err
	}

	resMap, ok := Result.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format from Midtrans")
	}

	jsonBytes, err := json.Marshal(resMap)
	if err != nil {
		return "", fmt.Errorf("failed to re-marshal result: %w", err)
	}

	if HttpStatus < 200 || HttpStatus >= 300 {
		var errRes MidtransErrorResponse
		if err := json.Unmarshal(jsonBytes, &errRes); err == nil && len(errRes.ErrorMessages) > 0 {
			return "", fmt.Errorf("midtrans error: %s", strings.Join(errRes.ErrorMessages, "; "))
		}
		return "", fmt.Errorf("midtrans returned HTTP %d but error message could not be parsed", HttpStatus)
	}

	var midtransRes MidtransSuccessResponse
	if err := json.Unmarshal(jsonBytes, &midtransRes); err != nil {
		return "", fmt.Errorf("failed to unmarshal to success struct: %w", err)
	}

	return midtransRes.RedirectURL, nil
}

func SendGetPaymentStatusToMidtrans(OrderID string, Vendor db_var.PaymentGatewayCredentialT) (string, string, error) {
	UrlEnvMode := global_var.PGUrlList.MidtransSend.Dev
	if Vendor.Mode == "prod" {
		UrlEnvMode = global_var.PGUrlList.MidtransSend.Prod
	}

	ApiKeys, err := helper.Decrypt(Vendor.APIKey, config.MasterKey)
	if err != nil {
		return "", "", err
	}

	Reqs := helper.RequestOptions{
		Method:   "GET",
		URL:      UrlEnvMode + "/v2/" + OrderID + "/status",
		AuthType: helper.AuthBasic,
		Username: ApiKeys,
	}

	Result, HttpStatus, _, err := helper.SendRequest(Reqs)
	if err != nil {
		return "", "", err
	}

	resMap, ok := Result.(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("unexpected response format from Midtrans")
	}

	jsonBytes, err := json.Marshal(resMap)
	if err != nil {
		return "", "", fmt.Errorf("failed to re-marshal result: %w", err)
	}

	if HttpStatus < 200 || HttpStatus >= 300 {
		var errRes MidtransErrorResponse
		if err := json.Unmarshal(jsonBytes, &errRes); err == nil && len(errRes.ErrorMessages) > 0 {
			return "", "", fmt.Errorf("midtrans error: %s", strings.Join(errRes.ErrorMessages, "; "))
		}
		return "", "", fmt.Errorf("midtrans returned HTTP %d but error message could not be parsed", HttpStatus)
	}

	var midtransRes MidtransNotificationStruct
	if err := json.Unmarshal(jsonBytes, &midtransRes); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal to success struct: %w", err)
	}

	return midtransRes.TransactionStatus, midtransRes.PaymentType, nil
}
