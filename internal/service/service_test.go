package service

import (
	"errors"
	"testing"
	"time"

	"bookkeeping/internal/config"
	"bookkeeping/internal/model"
	"bookkeeping/internal/store"
	"bookkeeping/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func (s *Service) seed(t *testing.T) (accountID, incomeCatID, expenseCatID string) {
	t.Helper()
	a, err := s.CreateAccount(model.Account{Name: "工资卡", Type: model.AccountBank, Balance: 100000})
	if err != nil {
		t.Fatalf("创建账户失败: %v", err)
	}
	inc, err := s.CreateCategory(model.Category{Name: "工资", Type: model.TypeIncome})
	if err != nil {
		t.Fatalf("创建收入分类失败: %v", err)
	}
	exp, err := s.CreateCategory(model.Category{Name: "餐饮", Type: model.TypeExpense})
	if err != nil {
		t.Fatalf("创建支出分类失败: %v", err)
	}
	return a.ID, inc.ID, exp.ID
}

func TestService_TransactionUpdatesBalance(t *testing.T) {
	s := newTestService()
	aid, incID, expID := s.seed(t)

	// 收入 5000 分
	if _, err := s.CreateTransaction(model.Transaction{AccountID: aid, CategoryID: incID, Type: model.TypeIncome, Amount: 5000}); err != nil {
		t.Fatalf("记账失败: %v", err)
	}
	// 支出 2000 分
	if _, err := s.CreateTransaction(model.Transaction{AccountID: aid, CategoryID: expID, Type: model.TypeExpense, Amount: 2000}); err != nil {
		t.Fatalf("记账失败: %v", err)
	}

	a, _ := s.GetAccount(aid)
	if a.Balance != 100000+5000-2000 {
		t.Fatalf("期望余额 %d，得到 %d", 100000+5000-2000, a.Balance)
	}
}

func TestService_ExpenseInsufficient(t *testing.T) {
	s := newTestService()
	aid, _, expID := s.seed(t)

	if _, err := s.CreateTransaction(model.Transaction{AccountID: aid, CategoryID: expID, Type: model.TypeExpense, Amount: 999999999}); !model.IsValidationError(err) {
		t.Fatalf("期望校验错误，得到 %v", err)
	}
}

func TestService_CategoryTypeMismatch(t *testing.T) {
	s := newTestService()
	aid, incID, _ := s.seed(t)

	// 收入分类却记支出
	if _, err := s.CreateTransaction(model.Transaction{AccountID: aid, CategoryID: incID, Type: model.TypeExpense, Amount: 100}); !model.IsValidationError(err) {
		t.Fatalf("期望校验错误，得到 %v", err)
	}
}

func TestService_DeleteTransactionRollback(t *testing.T) {
	s := newTestService()
	aid, _, expID := s.seed(t)

	tx, err := s.CreateTransaction(model.Transaction{AccountID: aid, CategoryID: expID, Type: model.TypeExpense, Amount: 3000})
	if err != nil {
		t.Fatalf("记账失败: %v", err)
	}
	before, _ := s.GetAccount(aid)
	beforeBalance := before.Balance

	if err := s.DeleteTransaction(tx.ID); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
	after, _ := s.GetAccount(aid)
	if after.Balance != beforeBalance+3000 {
		t.Fatalf("删除后余额应回滚，期望 %d，得到 %d", beforeBalance+3000, after.Balance)
	}
}

func TestService_MonthlyReportAndBudget(t *testing.T) {
	s := newTestService()
	aid, incID, expID := s.seed(t)

	now := time.Now()
	// 本月支出 8000 分
	if _, err := s.CreateTransaction(model.Transaction{AccountID: aid, CategoryID: expID, Type: model.TypeExpense, Amount: 8000, OccurredAt: now}); err != nil {
		t.Fatalf("记账失败: %v", err)
	}
	// 本月收入 10000 分
	if _, err := s.CreateTransaction(model.Transaction{AccountID: aid, CategoryID: incID, Type: model.TypeIncome, Amount: 10000, OccurredAt: now}); err != nil {
		t.Fatalf("记账失败: %v", err)
	}

	report, err := s.MonthlyReport(now.Year(), int(now.Month()))
	if err != nil {
		t.Fatalf("报表失败: %v", err)
	}
	if report.Income != 10000 || report.Expense != 8000 || report.Net != 2000 {
		t.Fatalf("报表异常: %+v", report)
	}

	// 预算 5000，本月支出 8000，应超支
	if _, err := s.CreateBudget(model.Budget{CategoryID: expID, Amount: 5000}); err != nil {
		t.Fatalf("创建预算失败: %v", err)
	}
	statuses, _ := s.CheckBudgets(now)
	if len(statuses) != 1 || !statuses[0].OverBudget {
		t.Fatalf("预算状态异常: %+v", statuses)
	}
}

func TestService_Overview(t *testing.T) {
	s := newTestService()
	s.seed(t)

	ov, err := s.Overview()
	if err != nil {
		t.Fatalf("总览失败: %v", err)
	}
	if ov.AccountCount != 1 || ov.TotalBalance != 100000 {
		t.Fatalf("总览异常: %+v", ov)
	}
	if _, err := s.MonthlyReport(0, 13); !errors.Is(err, store.ErrNotFound) && !model.IsValidationError(err) {
		t.Fatalf("非法月份应报错")
	}
}
