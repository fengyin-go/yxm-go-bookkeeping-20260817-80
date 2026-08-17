# Bug 是什么

分类已经被流水使用后，仍然允许删除分类，导致已有流水指向不存在的分类，月度报表中的分类名为空。

# 如何触发

创建分类和一笔流水后，删除该分类，再调用月度报表接口查看分类汇总。

# 错误信息

`go test ./internal/service -run TestService_DeleteCategoryUsedByTransactionIsRejected` 会失败，期望返回分类占用错误，但实际得到 `nil`。
