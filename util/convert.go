package util

import (
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
)

// Convert 将源对象按字段名复制到新建的目标类型实例(通用 mapper -> response 转换)
// 同名导出字段自动复制(类型可转换时自动转换, 如 int64 -> int);
// 字段名不同时在目标字段添加 copier:"源字段名" tag 映射;
// 源中无对应字段的目标字段保持零值; 基于 github.com/jinzhu/copier 反射实现
func Convert[T any](src any) (T, error) {
	var dst T
	err := copier.Copy(&dst, src)
	return dst, err
}

// ConvertSlice 将源对象切片按字段名复制到目标类型切片(逐元素同 Convert 规则)
// 始终返回非 nil 切片(源为空/nil 时返回空切片)
func ConvertSlice[T any](src any) ([]T, error) {
	dst := make([]T, 0)
	err := copier.Copy(&dst, src)
	return dst, err
}

func ToString(src any) string {
	return cast.ToString(src)
}
