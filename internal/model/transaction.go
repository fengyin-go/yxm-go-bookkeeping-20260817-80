package model

import (
	"strings"
	"time"
)

// Transaction 一笔收支流水。
type Transaction struct {
	ID         string    `json:"id"`
	AccountID  string    `json:"account_id"`
	CategoryID string    `json:"category_id"`
	Type       string    `json:"type"`       // income / expense
	Amount     int64     `json:"amount"`     // 金额（分，正数）
	Note       string    `json:"note"`       // 备注
	OccurredAt time.Time `json:"occurred_at"` // 发生时间
	CreatedAt  time.Time `json:"created_at"`
}

// Validate 规范化并校验流水字段。
func (t *Transaction) Validate() error {
	t.AccountID = strings.TrimSpace(t.AccountID)
	t.CategoryID = strings.TrimSpace(t.CategoryID)
	t.Note = strings.TrimSpace(t.Note)
	if t.AccountID == "" {
		return NewValidationError("account_id", "账户不能为空")
	}
	if t.CategoryID == "" {
		return NewValidationError("category_id", "分类不能为空")
	}
	if t.Type != TypeIncome && t.Type != TypeExpense {
		return NewValidationError("type", "收支类型需为 income 或 expense")
	}
	if t.Amount <= 0 {
		return NewValidationError("amount", "金额必须为正数")
	}
	if t.OccurredAt.IsZero() {
		t.OccurredAt = time.Now()
	}
	return nil
}

// TransactionFilter 流水筛选条件。
type TransactionFilter struct {
	AccountID  string
	CategoryID string
	Type       string
	From       *time.Time // 起始时间（含）
	To         *time.Time // 结束时间（含）
}

// Match 判断流水是否命中筛选条件。
func (f TransactionFilter) Match(t *Transaction) bool {
	if f.AccountID != "" && t.AccountID != f.AccountID {
		return false
	}
	if f.CategoryID != "" && t.CategoryID != f.CategoryID {
		return false
	}
	if f.Type != "" && t.Type != f.Type {
		return false
	}
	if f.From != nil && t.OccurredAt.Before(*f.From) {
		return false
	}
	if f.To != nil && t.OccurredAt.After(*f.To) {
		return false
	}
	return true
}
