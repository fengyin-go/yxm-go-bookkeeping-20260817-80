package model

import (
	"errors"
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

// CategoryInUseError 表示分类已经被流水使用，不能变更关键属性。
type CategoryInUseError struct {
	CategoryID string
}

func (e *CategoryInUseError) Error() string {
	if e.CategoryID != "" {
		return "分类已被流水使用: " + e.CategoryID
	}
	return "分类已被流水使用"
}

// IsCategoryInUse 判断错误是否为分类占用错误。
func IsCategoryInUse(err error) bool {
	var target *CategoryInUseError
	return errors.As(err, &target)
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
