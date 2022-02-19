package mathext

func Argmax(values []float64) (int, float64) {
	iMax, max := 0, values[0]
	for i, v := range values {
		if v > max {
			iMax, max = i, v
		}
	}
	return iMax, max
}
