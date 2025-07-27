package global_var

import (
	"gorm.io/gorm"
)

// Struct Section

type DatabaseConnection struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Host         string `json:"email"`
	Port         string `json:"port"`
	DatabaseName string `json:"database_name"`
}

type TRequestResponse struct {
	Message interface{} `json:"message"`
	Result  interface{} `json:"result"`
}

type TRequestMethod struct {
	Get, Post, Put, Patch, Delete, Update string
}

type PGVendorStruct struct {
	Midtrans, Xendit, HitPay, Doku string
}

type PGEnvStatus struct {
	Dev, Prod string
}

type PGEnvUrl struct {
	Midtrans     PGEnvStatus
	MidtransSend PGEnvStatus
	Xendit       PGEnvStatus
}

// Global Variable
var DB *gorm.DB
var RequestMethod = TRequestMethod{
	Post:   "POST",
	Get:    "GET",
	Put:    "PUT",
	Patch:  "PATCH",
	Delete: "DELETE",
	Update: "UPDATE"}

var PGVendor = PGVendorStruct{
	Midtrans: "MIDTR",
	Xendit:   "XNDT",
	HitPay:   "HTPY",
	Doku:     "DOKU",
}

var (
	TxStatusPending        = "pending"
	TxStatusSent           = "sent"
	TxStatusWaitingPayment = "waiting_payment"
	TxStatusPaid           = "paid"
	TxStatusExpired        = "expired"
	TxStatusFailed         = "failed"
	TxStatusError          = "error"
	TxStatusRefunded       = "refunded"
)

var PGUrlList = PGEnvUrl{
	Midtrans: PGEnvStatus{
		Dev:  "https://app.sandbox.midtrans.com",
		Prod: "https://app.midtrans.com",
	},
	MidtransSend: PGEnvStatus{
		Dev:  "https://api.sandbox.midtrans.com",
		Prod: "https://api.midtrans.com",
	},
	Xendit: PGEnvStatus{
		Dev:  "https://api.sandbox.xendit.co",
		Prod: "https://api.xendit.co",
	},
}
