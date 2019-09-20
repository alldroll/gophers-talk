package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const nGram = 2

var dictionary = []string{
	"gopher",
	"go",
	"london",
	"brother",
	"php",
	"philosophy",
	"independent",
}

func main() {
	suggester := BuildIndex(nGram, dictionary)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">>> ")

	for scanner.Scan() {
		query := strings.TrimSpace(scanner.Text())

		if len(query) == 0 {
			fmt.Print(">>> ")
			continue
		}

		result := suggester.Search(query)

		for _, item := range result {
			fmt.Printf("%s\n", item)
		}

		fmt.Print(">>> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Something went bad: %v", err)
	}
}
