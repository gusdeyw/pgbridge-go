package db_var

import (
	"time"

	"gorm.io/datatypes"
)

// Struct Section

type UserT struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"type:varchar(30);unique"`
	Password  string    `json:"password" gorm:"type:varchar(200)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (UserT) TableName() string {
	return TableName.User
}

type PaymentGatewayCredentialT struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Code             string    `json:"code" gorm:"type:varchar(100);uniqueIndex"`
	UserCode         string    `json:"user_code" gorm:"type:varchar(50);not null"`
	GatewayName      string    `json:"gateway_name" gorm:"type:varchar(50);not null"`
	APIKey           string    `json:"api_key" gorm:"type:varchar(200);not null"`
	APISecret        string    `json:"api_secret" gorm:"type:varchar(200);not null"`
	MerchantID       string    `json:"merchant_id" gorm:"type:varchar(100)"`
	CallbackURL      string    `json:"callback_url" gorm:"type:varchar(200)"`
	CallbackRedirect int       `json:"callback_redirect" gorm:"default:0"`
	Mode             string    `json:"mode" gorm:"type:varchar(10);default:'dev'"`
	CreatedAt        time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy        string    `json:"updated_by"`
}

func (PaymentGatewayCredentialT) TableName() string {
	return TableName.PGCredentials
}

type PaymentGatewayTransactionT struct {
	ID             uint64         `json:"id" gorm:"primaryKey"`
	OrderID        string         `json:"order_id" gorm:"type:varchar(64);uniqueIndex;not null"`
	UserCode       string         `json:"user_code" gorm:"type:varchar(50);not null"`
	Amount         int            `json:"amount" gorm:"not null"`
	CustomerName   string         `json:"customer_name" gorm:"type:varchar(255)"`
	CustomerEmail  string         `json:"customer_email" gorm:"type:varchar(255)"`
	CustomerPhone  string         `json:"customer_phone" gorm:"type:varchar(50)"`
	ItemsJSON      datatypes.JSON `json:"items_json" gorm:"type:jsonb"`
	PaymentMethods string         `json:"payment_methods"`
	CustomFields   datatypes.JSON `json:"custom_fields" gorm:"type:jsonb"`
	Metadata       datatypes.JSON `json:"metadata" gorm:"type:jsonb"`
	CallbacksJSON  datatypes.JSON `json:"callbacks_json" gorm:"type:jsonb"`
	ExpiryStart    *time.Time     `json:"expiry_start"`
	ExpiryUnit     string         `json:"expiry_unit" gorm:"type:varchar(20)"`
	ExpiryDuration int            `json:"expiry_duration"`
	Vendor         string         `json:"vendor" gorm:"type:varchar(50)"`
	VendorPayload  datatypes.JSON `json:"vendor_payload" gorm:"type:jsonb"`
	Status         string         `json:"status" gorm:"type:varchar(50);default:'pending'"`
	PaidAt         time.Time      `json:"paid_at"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy string    `json:"updated_by"`
}

func (PaymentGatewayTransactionT) TableName() string {
	return TableName.PGTransactions // use your constants package
}

// Variable

// list of table name
type TableNameStruct struct {
	User           string
	PGCredentials  string
	PGTransactions string
}

var TableName = TableNameStruct{
	User:           "user",
	PGCredentials:  "payment_gateway_credentials",
	PGTransactions: "payment_gateway_transaction",
}
