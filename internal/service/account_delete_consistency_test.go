package service

import (
	"testing"

	"bookkeeping/internal/model"
)

func TestService_DeleteAccountUsedByTransactionIsRejected(t *testing.T) {
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

	err := s.DeleteAccount(aid)
	if !model.IsAccountInUse(err) {
		t.Fatalf("期望账户占用错误，得到 %v", err)
	}
}
