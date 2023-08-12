package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	input, err := os.Open("input.txt")
	defer input.Close()
	if err != nil {
		log.Fatal(err)
	}
	max := 0
	currentTotal := 0
	sc := bufio.NewScanner(input)
	for sc.Scan() {
		currentCalories, err := strconv.Atoi(sc.Text())
		if err != nil {
			//Means current cals for elf has been found
			if max < currentTotal {
				max = currentTotal
			}
			currentTotal = 0
            continue
		}
		currentTotal += currentCalories
	}
	fmt.Println(max)

}
