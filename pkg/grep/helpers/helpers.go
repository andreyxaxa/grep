package helpers

import (
	"bufio"
	"io"
	"log"
)

// ReadLines reads from 'from' to 'to' (bufio.Scanner)
func ReadLines(from io.Reader, to *[]string) error {
	scanner := bufio.NewScanner(from)

	for scanner.Scan() {
		*to = append(*to, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error while scan: %v\n", err)
	}

	return nil
}
