package main

func combinations(pool []TokenPointer, r int) [][]TokenPointer {
	lenpool := len(pool)
	result := [][]TokenPointer{}
	if r > lenpool {
		return result
	}

	combination := []TokenPointer{}
	indices := []int{}
	for i := 0; i < r; i++ {
		combination = append(combination, pool[i])
		indices = append(indices, i)
	}
	result = append(result, combination)

	for true {
		broke := false
		var i int
		for i = r - 1; i >= 0; i-- {
			if indices[i] != i+lenpool-r {
				broke = true
				break
			}
		}
		if !broke {
			break
		}
		indices[i] += 1
		for j := i + 1; j < r; j++ {
			indices[j] = indices[j-1] + 1
		}

		combination = []TokenPointer{}
		for _, ind := range indices {
			combination = append(combination, pool[ind])
		}
		result = append(result, combination)
	}
	return result
}
