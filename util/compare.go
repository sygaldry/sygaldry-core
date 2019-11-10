package util

// ListContainsListStrings checks if large contains all elements of small list
func ListContainsListStrings(smallList []string, largeList []string) bool {
	for _, element := range smallList {
		if IndexOfElementStrings(largeList, element) < 0 && element != "definition" {
			return false
		}
	}
	return true
}

// IndexOfElementStrings checks if list contains element
func IndexOfElementStrings(list []string, element string) int {
	for index, listElement := range list {
		if listElement == element {
			return index
		}
	}
	return -1
}
