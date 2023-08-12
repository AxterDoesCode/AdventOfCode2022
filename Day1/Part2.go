package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	input, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()
	elves := []int{}
	sc := bufio.NewScanner(input)
	currentCalTotal := 0
	for sc.Scan() {
		currentCals, err := strconv.Atoi(sc.Text())
		if err != nil {
			elves = append(elves, currentCalTotal)
			currentCalTotal = 0
			continue
		}
		currentCalTotal += currentCals
	}
    sort.Slice(elves, func(i, j int) bool {
        return elves[i] > elves[j]
    })
    fmt.Println(elves[0] + elves[1] + elves[2])
}
