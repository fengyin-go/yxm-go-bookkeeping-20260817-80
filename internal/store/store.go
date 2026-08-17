// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"bookkeeping/internal/model"
)

var (
	ErrNotFound            = errors.New("记录不存在")
	ErrConflict            = errors.New("记录已存在或状态冲突")
	ErrInsufficientBalance = errors.New("账户余额不足")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	// 账户
	CreateAccount(a *model.Account) error
	GetAccount(id string) (*model.Account, error)
	GetAccountByName(name string) (*model.Account, error)
	ListAccounts() []*model.Account
	UpdateAccount(a *model.Account) error
	AdjustAccountBalance(id string, delta int64) (*model.Account, error)
	DeleteAccount(id string) error

	// 分类
	CreateCategory(c *model.Category) error
	GetCategory(id string) (*model.Category, error)
	ListCategories() []*model.Category
	UpdateCategory(c *model.Category) error
	DeleteCategory(id string) error

	// 流水
	CreateTransaction(t *model.Transaction) error
	GetTransaction(id string) (*model.Transaction, error)
	ListTransactions() []*model.Transaction
	DeleteTransaction(id string) error

	// 预算
	CreateBudget(b *model.Budget) error
	GetBudget(id string) (*model.Budget, error)
	ListBudgets() []*model.Budget
	UpdateBudget(b *model.Budget) error
	DeleteBudget(id string) error
}
