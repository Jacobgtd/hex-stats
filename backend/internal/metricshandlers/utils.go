package metricshandlers

func absInt(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
