package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main () {
    input, err := os.Open("input.txt")
     
    if err != nil {
        log.Fatal(err)
        return
    }
    sc := bufio.NewScanner(input)

    var score int

    // A : X : Rock
    // B : Y : Paper
    // C : Z : Scissors

    scores := map[string]int{"B X" : 1, "C Y" : 2, "A Z" : 3, "A X" : 4, "B Y" : 5, "C Z" : 6, "C X" : 7, "A Y" : 8, "B Z" : 9}
    
    for sc.Scan() {
        score += scores[sc.Text()]
    }
    fmt.Println(score)
}
