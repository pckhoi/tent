package utils

func MergeStringSlices(dstSlice []string, srcSlice []string) []string {
	newSlice := make([]string, len(dstSlice))
	copy(newSlice, dstSlice)
	for _, srcString := range srcSlice {
		found := false
		for _, dstString := range dstSlice {
			if srcString == dstString {
				found = true
				break
			}
		}
		if !found {
			newSlice = append(newSlice, srcString)
		}
	}
	return newSlice
}

func StringSliceEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
