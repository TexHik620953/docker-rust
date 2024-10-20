package utils

func RefValue[T any](value T) *T {
	t := value
	return &t
}
