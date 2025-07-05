package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <file1.txt> <file2.txt> <output.txt>")
		fmt.Println("Example: go run main.go data/disposable.txt data/free.txt merged_domains.txt")
		os.Exit(1)
	}

	file1 := os.Args[1]
	file2 := os.Args[2]
	outputFile := os.Args[3]

	// Read domains from both files
	domains := make(map[string]bool)

	// Read first file
	if err := readDomainsFromFile(file1, domains); err != nil {
		fmt.Printf("Error reading %s: %v\n", file1, err)
		os.Exit(1)
	}

	// Read second file
	if err := readDomainsFromFile(file2, domains); err != nil {
		fmt.Printf("Error reading %s: %v\n", file2, err)
		os.Exit(1)
	}

	// Convert map keys to slice for sorting
	var domainList []string
	for domain := range domains {
		domainList = append(domainList, domain)
	}

	// Sort domains
	sort.Strings(domainList)

	// Write sorted domains to output file
	if err := writeDomainsToFile(outputFile, domainList); err != nil {
		fmt.Printf("Error writing to %s: %v\n", outputFile, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully merged %s and %s into %s\n", file1, file2, outputFile)
	fmt.Printf("Total unique domains: %d\n", len(domainList))
}

func readDomainsFromFile(filename string, domains map[string]bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			domains[strings.ToLower(domain)] = true
		}
	}

	return scanner.Err()
}

func writeDomainsToFile(filename string, domains []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, domain := range domains {
		_, err := writer.WriteString(domain + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
