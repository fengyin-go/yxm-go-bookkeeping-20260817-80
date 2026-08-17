package service

import (
	"testing"

	"bookkeeping/internal/model"
)

func TestService_UpdateCategoryTypeUsedByTransactionIsRejected(t *testing.T) {
	s := newTestService()
	aid, _, expID := s.seed(t)

	if _, err := s.CreateTransaction(model.Transaction{
		AccountID:  aid,
		CategoryID: expID,
		Type:       model.TypeExpense,
		Amount:     100,
	}); err != nil {
		t.Fatalf("创建支出流水失败: %v", err)
	}

	_, err := s.UpdateCategory(expID, model.Category{Name: "餐饮", Type: model.TypeIncome})
	if !model.IsCategoryInUse(err) {
		t.Fatalf("期望分类占用错误，得到 %v", err)
	}
}

// 分类被流水占用时，只改名称（类型不变）应当允许。
func TestService_UpdateCategoryNameOnlyWhenInUse(t *testing.T) {
	s := newTestService()
	aid, _, expID := s.seed(t)

	if _, err := s.CreateTransaction(model.Transaction{
		AccountID:  aid,
		CategoryID: expID,
		Type:       model.TypeExpense,
		Amount:     100,
	}); err != nil {
		t.Fatalf("创建支出流水失败: %v", err)
	}

	updated, err := s.UpdateCategory(expID, model.Category{Name: "吃饭", Type: model.TypeExpense})
	if err != nil {
		t.Fatalf("仅改名称不应报错，得到 %v", err)
	}
	if updated.Name != "吃饭" || updated.Type != model.TypeExpense {
		t.Fatalf("分类更新异常: %+v", updated)
	}
}

// 分类未被流水占用时，变更类型应当允许。
func TestService_UpdateCategoryTypeWhenNotInUse(t *testing.T) {
	s := newTestService()
	_, _, expID := s.seed(t)

	updated, err := s.UpdateCategory(expID, model.Category{Name: "餐饮", Type: model.TypeIncome})
	if err != nil {
		t.Fatalf("未被占用时变更类型不应报错，得到 %v", err)
	}
	if updated.Type != model.TypeIncome {
		t.Fatalf("类型应已更新为 income: %+v", updated)
	}
}
