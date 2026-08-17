package model

import (
	"strings"
	"time"
)

// 账户类型常量。
const (
	AccountCash   = "cash"   // 现金
	AccountBank   = "bank"   // 银行卡
	AccountCredit = "credit" // 信用卡
	AccountOther  = "other"  // 其他
)

// Account 资金账户。
type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`     // cash / bank / credit / other
	Balance   int64     `json:"balance"`  // 余额（分）
	Currency  string    `json:"currency"` // 币种，默认 CNY
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 规范化并校验账户字段。
func (a *Account) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	if a.Name == "" {
		return NewValidationError("name", "账户名称不能为空")
	}
	if a.Type == "" {
		a.Type = AccountCash
	}
	if !validAccountType(a.Type) {
		return NewValidationError("type", "账户类型不合法")
	}
	if a.Balance < 0 {
		return NewValidationError("balance", "初始余额不能为负数")
	}
	if a.Currency == "" {
		a.Currency = "CNY"
	}
	return nil
}

// validAccountType 校验账户类型。
func validAccountType(t string) bool {
	switch t {
	case AccountCash, AccountBank, AccountCredit, AccountOther:
		return true
	default:
		return false
	}
}

// ValidAccountType 对外暴露的账户类型校验。
func ValidAccountType(t string) bool { return validAccountType(t) }
