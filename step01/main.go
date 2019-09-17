package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	dictionary := []string{
		"gopher",
		"go",
		"london",
		"brother",
		"php",
		"philosophy",
		"independent",
	}

	suggester := BuildIndex(2, dictionary)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")

	for scanner.Scan() {
		query := strings.TrimSpace(scanner.Text())

		if len(query) == 0 {
			fmt.Print(">> ")
			continue
		}

		result := suggester.Search(query)

		for _, item := range result {
			fmt.Printf("%s\n", item)
		}

		fmt.Print(">> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Errorf("Something went bad: %v", err)
	}
}
