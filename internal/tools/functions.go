package tools

func StringIsInArray(element string, array []string) bool {
	for _, item := range array {
		if item == element {
			return true
		}
	}
	return false
}
