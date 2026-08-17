package model

import (
	"errors"
	"strings"
	"time"
)

// 预算周期常量。
const (
	PeriodMonthly = "monthly"
)

// Budget 预算，用于监控某分类（或全部支出）的月度限额。
type Budget struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"category_id"` // 为空表示全部支出
	Period     string    `json:"period"`      // 固定 monthly
	Amount     int64     `json:"amount"`      // 预算金额（分）
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BudgetDuplicateError 表示同一分类（或全局）已经存在预算。
type BudgetDuplicateError struct {
	CategoryID string
}

func (e *BudgetDuplicateError) Error() string {
	if e.CategoryID != "" {
		return "该分类已存在预算: " + e.CategoryID
	}
	return "已存在全局预算"
}

// IsBudgetDuplicate 判断错误是否为重复预算错误。
func IsBudgetDuplicate(err error) bool {
	var target *BudgetDuplicateError
	return errors.As(err, &target)
}

// Validate 规范化并校验预算字段。
func (b *Budget) Validate() error {
	b.CategoryID = strings.TrimSpace(b.CategoryID)
	if b.Amount <= 0 {
		return NewValidationError("amount", "预算金额必须为正数")
	}
	if b.Period == "" {
		b.Period = PeriodMonthly
	}
	if b.Period != PeriodMonthly {
		return NewValidationError("period", "预算周期需为 monthly")
	}
	return nil
}
