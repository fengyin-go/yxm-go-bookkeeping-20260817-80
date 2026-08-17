package model

import (
	"strings"
	"time"
)

// 收支类型常量。
const (
	TypeIncome  = "income"  // 收入
	TypeExpense = "expense" // 支出
)

// Category 收支分类。
type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // income / expense
	CreatedAt time.Time `json:"created_at"`
}

// Validate 规范化并校验分类字段。
func (c *Category) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return NewValidationError("name", "分类名称不能为空")
	}
	if c.Type == "" {
		return NewValidationError("type", "分类类型不能为空")
	}
	if c.Type != TypeIncome && c.Type != TypeExpense {
		return NewValidationError("type", "分类类型需为 income 或 expense")
	}
	return nil
}

// ValidType 校验收支类型。
func ValidType(t string) bool {
	return t == TypeIncome || t == TypeExpense
}
