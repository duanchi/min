package util

// Nilable 安全解引用指针: nil 返回零值
func Nilable[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
