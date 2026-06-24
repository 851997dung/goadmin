package buildin

func If(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

func IfInt(b bool, t, f int) int {
	if b {
		return t
	}
	return f
}
