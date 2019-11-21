package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const nGram = 2
const topK = 10

func readDictionary() ([]string, error) {
	file, err := os.Open("cities.dict")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var dictionary []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		dictionary = append(dictionary, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dictionary, nil
}

func main() {
	dictionary, err := readDictionary()

	if err != nil {
		log.Fatalf("Failed to read dictionary: %v", err)
	}

	suggester := BuildIndex(nGram, dictionary)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">>> ")

	for scanner.Scan() {
		query := strings.TrimSpace(scanner.Text())

		if len(query) == 0 {
			fmt.Print(">>> ")
			continue
		}

		start := time.Now()
		result := suggester.Search(query, topK)
		elapsed := time.Since(start).String()

		for _, item := range result {
			fmt.Printf("%s: %f\n", item.Candidate, item.Score)
		}

		fmt.Printf("\nElapsed: %s (%d candidates)\n", elapsed, len(result))
		fmt.Print(">>> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Something went bad: %v", err)
	}
}
