// The main package of the concurrent matrix determinant calculation
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

// RandomLimit - limit of numbers in the generated matrix
const RandomLimit int = 100
const noFile string = "-"

// get console arguments and start program
func main() {
	var t, n int
	var i string
	flag.IntVar(&t, "t", 3, "The number of goroutines to be started")
	flag.IntVar(&n, "n", 3, "The size of the matrix to be generated")
	flag.StringVar(&i, "i", noFile, "Read the matrix from a file (overrides n)")
	flag.Parse()
	var matrix [][]int
	var err error
	if i != noFile {
		matrix, err = readMatrixFromFile(i)
		if err != nil {
			return
		}
	} else {
		matrix = make([][]int, n)
		for i := 0; i < n; i++ {
			matrix[i] = make([]int, n)
		}
		randomizeMatrix(matrix)
	}
	indexes := make([]int, n)
	for i := 0; i < n; i++ {
		indexes[i] = i
	}
	permutations := getPermutations(indexes, n)
	start := time.Now()
	detChannels := make([]<-chan int, t)
	piece := len(permutations) / t
	for i := 1; i <= t; i++ {
		detChannels[i-1] = determinant(matrix, permutations[(i-1)*piece:i*piece])
	}
	fmt.Println("D =", sumChannels(detChannels))
	elapsed := time.Since(start)
	fmt.Printf(
		"Calculation took %fs for n=%d and %d routine(s)\n",
		elapsed.Seconds(),
		len(matrix[0]),
		t)
}

func readMatrixFromFile(fileName string) ([][]int, error) {
	var matrix [][]int
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	n := 0
	nI := 0
	for _, symbol := range file {
		if symbol == '\n' {
			break
		}
		n *= 10
		dig, _ := strconv.Atoi(string(symbol))
		n += dig
		nI++
	}
	matrix = make([][]int, n)
	for k := 0; k < n; k++ {
		matrix[k] = make([]int, n)
	}
	i, j := 0, 0
	currentNumber := 0
	for _, symbol := range file[nI+1:] {
		if symbol >= '0' && symbol <= '9' {
			currentNumber *= 10
			dig, _ := strconv.Atoi(string(symbol))
			currentNumber += dig
		} else if symbol == ' ' {
			matrix[i][j] = currentNumber
			j++
			currentNumber = 0
		} else if symbol == '\n' {
			matrix[i][j] = currentNumber
			i++
			j = 0
			currentNumber = 0
		}
	}
	return matrix, nil
}

// It's a kind of magic
func sumChannels(channels []<-chan int) (sum int) {
	aggregate := make(chan int)
	for _, c := range channels {
		go func(c <-chan int) {
			aggregate <- <-c
		}(c)
	}
	for _ = range channels {
		addend := <-aggregate
		sum += addend
	}
	return
}

func determinant(matrix [][]int, permutations [][]int) <-chan int {
	det := make(chan int)
	go func() {
		detNum := 0
		n := len(matrix)
		rowIndexes := make([]int, n)
		for i := 0; i < n; i++ {
			rowIndexes[i] = i
		}
		for _, permutation := range permutations {
			sign := getDeterminantMultipleSign(permutation)
			multiple := 1
			for i, cIndex := range permutation {
				multiple *= matrix[i][cIndex]
			}
			detNum += sign * multiple
		}
		det <- detNum
	}()
	return det
}

func randomizeMatrix(matrix [][]int) {
	n := len(matrix)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			matrix[i][j] = r1.Intn(RandomLimit)
		}
	}
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

func getDeterminantMultipleSign(permutation []int) (multiplier int) {
	n := len(permutation)
	inversions := 0
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if permutation[i] > permutation[j] {
				inversions++
			}
		}
	}
	if inversions%2 == 0 {
		return 1
	}
	return -1
}

// permute calculates the permutations of numbers in the permutations slice
// TODO verify and make concurrent
func getPermutations(iterable []int, r int) (permutations [][]int) {
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
