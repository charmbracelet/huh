package huh

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp(n, low, high int) int {
	if low > high {
		low, high = high, low
	}
	return min(high, max(low, n))
}
