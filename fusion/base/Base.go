package base

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func CopyBytes(src []byte) []byte {
	if src != nil {
		dst := make([]byte, len(src))
		copy(dst, src)
		return dst
	}
	return nil
}
