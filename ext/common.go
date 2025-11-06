package ext

// ternary operator approximation.
func iftrue[T any](b bool, t T, f T) T {
	if b {
		return t
	}
	return f
}
