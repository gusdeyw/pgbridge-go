package global_var

import "time"

// struct

type XDNT_ItemDetail_RequestBody struct {
	ReferenceID string  `json:"reference_id"`
	Name        string  `json:"name"`
	Currency    string  `json:"currency"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
}

type XDNT_RequestBody struct {
	ReferenceID string                        `json:"reference_id"`
	Type        string                        `json:"type"`
	Currency    string                        `json:"currency"`
	Amount      float64                       `json:"amount"`
	ExpiresAt   time.Time                     `json:"expires_at"`
	Basket      []XDNT_ItemDetail_RequestBody `json:"basket"`
}

type XDNT_ItemDetail_ResultBody struct {
	ReferenceID string  `json:"reference_id"`
	Name        string  `json:"name"`
	Currency    string  `json:"currency"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
}

type XDNT_ResultBody struct {
	ReferenceID string                       `json:"reference_id"`
	Type        string                       `json:"type"`
	Currency    string                       `json:"currency"`
	ChannelCode string                       `json:"channel_code"`
	Amount      float64                      `json:"amount"`
	ExpiresAt   time.Time                    `json:"expires_at"`
	Basket      []XDNT_ItemDetail_ResultBody `json:"basket"`
	BusinessID  string                       `json:"business_id"`
	ID          string                       `json:"id"`
	Created     time.Time                    `json:"created"`
	Updated     time.Time                    `json:"updated"`
	QRString    string                       `json:"qr_string"`
	Status      string                       `json:"status"`
}
