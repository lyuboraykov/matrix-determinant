// The main package of the concurrent matrix determinant calculation
package main

import "fmt"

// get console arguments and start program
func main() {
	// TODO: get n from console input; n - matrix size
	var n = 4
	indexes := make([]int, n)
	for i := 0; i < n; i++ {
		indexes[i] = i
	}
	permutations := permutate(indexes, n)
	fmt.Println(permutations)
}

// fact calculates the n factorial iteratively
// TODO: make concurrent
func fact(n int) (fact int) {
	fact = 1
	for i := 1; i <= n; i++ {
		fact *= i
	}
	return
}

// permute calculates the permutations of numbers in the permutations slice
// TODO verify and make concurrent
func permutate(iterable []int, r int) (permutations [][]int) {
	n := len(iterable)
	nFact := fact(n)
	for i := 0; i < nFact; i++ {
		permutations = append(permutations, make([]int, n))
	}
	pool := iterable

	if r > n {
		return
	}

	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	cycles := make([]int, r)
	for i := range cycles {
		cycles[i] = n - i
	}

	result := make([]int, r)
	for i, el := range indices[:r] {
		result[i] = pool[el]
	}
	index := 0
	copy(permutations[index], result)
	index++

	for n > 0 {
		i := r - 1
		for ; i >= 0; i-- {
			cycles[i]--
			if cycles[i] == 0 {
				index := indices[i]
				for j := i; j < n-1; j++ {
					indices[j] = indices[j+1]
				}
				indices[n-1] = index
				cycles[i] = n - i
			} else {
				j := cycles[i]
				indices[i], indices[n-j] = indices[n-j], indices[i]

				for k := i; k < r; k++ {
					result[k] = pool[indices[k]]
				}
				copy(permutations[index], result)
				index++
				break
			}
		}
		if i < 0 {
			return
		}
	}
	return
}
