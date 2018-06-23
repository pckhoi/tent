package main

import (
	"strconv"
)

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
		indices[i]++
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

func nameWithCount() func(string) string {
	counter := map[string]int{}

	return func(str string) string {
		_, ok := counter[str]
		if ok {
			counter[str]++
		} else {
			counter[str] = 1
		}
		return str + strconv.Itoa(counter[str])
	}
}
