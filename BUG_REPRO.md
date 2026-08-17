# Bug 是什么

账户已经被流水使用后，仍然允许直接删除账户，导致流水成为孤儿，余额也无法再回滚。

# 如何触发

创建账户和一笔支出流水后，调用删除账户接口，再查询流水。

# 错误信息

`go test ./internal/service -run TestService_DeleteAccountUsedByTransactionIsRejected` 会失败，期望返回账户占用错误，但实际得到 `nil`。
