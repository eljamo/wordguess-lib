package validator

func IsElementInSlice[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func UnorderedSlicesAreEqual[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	elementCount := make(map[T]int)

	// Count the frequency of each element in slice1
	for _, elem := range slice1 {
		elementCount[elem]++
	}

	// Subtract the frequency of each element in slice2
	for _, elem := range slice2 {
		if _, ok := elementCount[elem]; !ok {
			return false // Element in slice2 not found in slice1
		}
		elementCount[elem]--
		if elementCount[elem] == 0 {
			delete(elementCount, elem)
		}
	}

	// If elementCount is empty, the slices are equal
	return len(elementCount) == 0
}
