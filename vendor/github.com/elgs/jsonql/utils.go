package jsonql

// CompareSlices - compares to slices, return false if they are not equal, true if they are.
func CompareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// ReverseString - reverses a string.
func ReverseString(input string) string {
	r := []rune(input)
	for i := 0; i < len(r)/2; i++ {
		j := len(r) - i - 1
		r[i], r[j] = r[j], r[i]
	}
	output := string(r)
	return output
}
