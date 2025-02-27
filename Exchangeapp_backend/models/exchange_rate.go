package models

import "time"

type ExchangeRate struct {
	// gorm:"primary_key"：表示该字段是数据库表的主键（Primary Key）
	ID uint `gorm:"primary_key" json:"_id"`
	// json:"from_currency"：表示在 JSON 中该字段的名称为 from_currency
	// binding:"required"：表示该字段是必填项，通常用于数据验证（如 Gin 框架中的绑定验证）
	FromCurrency string    `json:"from_currency" binding:"required"`
	ToCurrency   string    `json:"to_currency" binding:"required"`
	Rate         float64   `json:"rate" binding:"required"`
	Date         time.Time `json:"date"`
}
