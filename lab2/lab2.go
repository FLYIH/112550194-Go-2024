package main

import "fmt"

func main() {
	var n int64

	fmt.Print("Enter a number: ")
	fmt.Scanln(&n)

	result := Sum(n)
	fmt.Println(result)
}

func Sum(n int64) string {
	// TODO: Finish this function
	var sum int64
	var steps string

	for i := int64(1); i <= n; i++ {
		if i%7 != 0 {
			sum += i
			if steps == "" {
				steps = fmt.Sprintf("%d", i)
			} else {
				steps += fmt.Sprintf("+%d", i)
			}
		}
	}
	steps += fmt.Sprintf("=%d", sum)
	return steps
}