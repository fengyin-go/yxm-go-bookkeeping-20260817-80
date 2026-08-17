# Bug 是什么

分类已经被流水使用后，仍然允许把分类从 `expense` 改成 `income`，破坏流水与分类类型的一致性。

# 如何触发

先创建一个支出分类，再创建一条该分类的支出流水，然后调用更新分类接口把类型改成 `income`。

# 错误信息

`go test ./internal/service -run TestService_UpdateCategoryTypeUsedByTransactionIsRejected` 会失败，期望返回分类占用错误，但实际得到 `nil`。
