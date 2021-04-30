package util

func StringsContains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func StringsContainsAll(slice []string, elements []string) bool {
	for _, e := range elements {
		if !StringsContains(slice, e) {
			return false
		}
	}
	return true
}
