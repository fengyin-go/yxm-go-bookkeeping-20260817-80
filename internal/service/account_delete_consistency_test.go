package service

import (
	"errors"
	"testing"

	"bookkeeping/internal/model"
	"bookkeeping/internal/store"
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

	// 账户与流水都应保留，余额不变。
	if _, err := s.GetAccount(aid); err != nil {
		t.Fatalf("拒绝删除后账户应仍存在: %v", err)
	}
	a, _ := s.GetAccount(aid)
	if a.Balance != 100000-100 {
		t.Fatalf("余额应保持一致，期望 %d，得到 %d", 100000-100, a.Balance)
	}
}

func TestService_DeleteAccountWithoutTransactionsSucceeds(t *testing.T) {
	s := newTestService()
	aid, _, _ := s.seed(t)

	// 账户无流水，应允许删除。
	if err := s.DeleteAccount(aid); err != nil {
		t.Fatalf("无流水的账户应可删除，得到 %v", err)
	}
	if _, err := s.GetAccount(aid); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("删除后账户应不存在，得到 %v", err)
	}
}

func TestService_DeleteNonExistentAccountReturnsNotFound(t *testing.T) {
	s := newTestService()
	if err := s.DeleteAccount("no-such-account"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
}
