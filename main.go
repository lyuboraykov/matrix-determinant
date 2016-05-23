// The main package of the concurrent matrix determinant calculation
package main

import "fmt"

// get console arguments and start program
func main() {
	// TODO: get n from console input; n - matrix size
	var n = 10
	indexes := make([]int, n)
	for i := 0; i < n; i++ {
		indexes[i] = i
	}
	nFact := fact(n)
	permutations := [][]int{}
	for i := 0; i < nFact; i++ {
		permutations = append(permutations, make([]int, n))
	}
	permute(indexes, permutations)
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
func permute(numbers []int, permutations [][]int) {
	if len(numbers) == 0 {
		return
	}
	currentI := len(permutations[0]) - len(numbers)
	for i, number := range numbers {
		permutations[i][currentI] = number
		permute(append(numbers[:currentI], numbers[currentI+1:]...), permutations)
	}
}
