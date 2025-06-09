package utilities

func To[T any](a T) *T {
	return &a
}
