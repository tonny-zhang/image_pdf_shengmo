package utils

// AbsUint8 get a-b abs
func AbsUint8(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}
