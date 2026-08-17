package service

import (
	"testing"

	"bookkeeping/internal/model"
)

func TestService_DuplicateBudgetIsRejected(t *testing.T) {
	s := newTestService()
	_, _, expID := s.seed(t)

	if _, err := s.CreateBudget(model.Budget{CategoryID: expID, Amount: 1000}); err != nil {
		t.Fatalf("首次创建预算失败: %v", err)
	}
	_, err := s.CreateBudget(model.Budget{CategoryID: expID, Amount: 2000})
	if !model.IsBudgetDuplicate(err) {
		t.Fatalf("期望重复预算错误，得到 %v", err)
	}
}
