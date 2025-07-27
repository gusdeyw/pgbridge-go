package models

import (
	"pg_bridge_go/db_var"
	"pg_bridge_go/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SavePGTransaction(db *gorm.DB, tx *db_var.PaymentGatewayTransactionT) error {
	if err := db.Create(tx).Error; err != nil {
		logger.Error("Failed to save transaction", zap.Error(err))
		return err
	}
	return nil
}

func InsertPGTransaction(txData *db_var.PaymentGatewayTransactionT, tx *gorm.DB) error {
	result := tx.Create(txData)
	return result.Error
}

func UpdatePGTransactionStatus(orderID, newStatus, PaymentMethods, updatedBy string, tx *gorm.DB) error {
	result := tx.
		Model(&db_var.PaymentGatewayTransactionT{}).
		Where("order_id = ?", orderID).
		Updates(map[string]interface{}{
			"payment_methods": PaymentMethods,
			"status":          newStatus,
			"updated_by":      updatedBy,
			"paid_at":         time.Now(),
			"updated_at":      time.Now(),
		})

	return result.Error
}
