package utils

func StringSliceExcept(slice []string, except []string) []string {
	newSlice := []string{}

	for _, c := range slice {
		if pos := strAt(except, c); pos == -1 {
			newSlice = append(newSlice, c)
		}
	}
	return newSlice
}

func StringSliceContainAll(container []string, slice []string) bool {
	if len(container) < len(slice) {
		return false
	}
	for _, s := range slice {
		if pos := strAt(container, s); pos == -1 {
			return false
		}
	}
	return true
}

func strAt(slice []string, str string) int {
	for pos, s := range slice {
		if s == str {
			return pos
		}
	}
	return -1
}
