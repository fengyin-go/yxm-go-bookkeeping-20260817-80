package store

import (
	"testing"
	"time"

	"bookkeeping/internal/model"
)

func TestMemoryStore_Account(t *testing.T) {
	s := NewMemoryStore()
	a := &model.Account{ID: "a1", Name: "工资卡", Type: model.AccountBank, Balance: 100000, Currency: "CNY", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateAccount(a); err != nil {
		t.Fatalf("创建账户失败: %v", err)
	}
	if err := s.CreateAccount(&model.Account{ID: "a2", Name: "工资卡"}); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
	got, err := s.GetAccountByName("工资卡")
	if err != nil || got.ID != "a1" {
		t.Fatalf("按名称查询失败: %v, %v", got, err)
	}
	a.Balance = 50000
	if err := s.UpdateAccount(a); err != nil {
		t.Fatalf("更新失败: %v", err)
	}
	after, _ := s.GetAccount("a1")
	if after.Balance != 50000 {
		t.Fatalf("期望余额 50000，得到 %d", after.Balance)
	}
}

func TestMemoryStore_Category(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Category{ID: "c1", Name: "餐饮", Type: model.TypeExpense, CreatedAt: time.Now()}
	if err := s.CreateCategory(c); err != nil {
		t.Fatalf("创建分类失败: %v", err)
	}
	// 同名不同类型允许
	if err := s.CreateCategory(&model.Category{ID: "c2", Name: "餐饮", Type: model.TypeIncome}); err != nil {
		t.Fatalf("同名不同类型应允许: %v", err)
	}
	// 同名同类型冲突
	if err := s.CreateCategory(&model.Category{ID: "c3", Name: "餐饮", Type: model.TypeExpense}); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
}

func TestMemoryStore_TransactionAndBudget(t *testing.T) {
	s := NewMemoryStore()
	tx := &model.Transaction{ID: "t1", AccountID: "a1", CategoryID: "c1", Type: model.TypeExpense, Amount: 1000, OccurredAt: time.Now(), CreatedAt: time.Now()}
	if err := s.CreateTransaction(tx); err != nil {
		t.Fatalf("创建流水失败: %v", err)
	}
	if len(s.ListTransactions()) != 1 {
		t.Fatalf("期望 1 条流水")
	}

	b := &model.Budget{ID: "b1", CategoryID: "c1", Amount: 5000, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateBudget(b); err != nil {
		t.Fatalf("创建预算失败: %v", err)
	}
	if len(s.ListBudgets()) != 1 {
		t.Fatalf("期望 1 条预算")
	}
	if err := s.DeleteTransaction("t1"); err != nil {
		t.Fatalf("删除流水失败: %v", err)
	}
}
