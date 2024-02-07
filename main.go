package main

import (
	"fmt"
	"os"
)

type segment struct {
	offset int
	length int
	name   string
	bytes  []byte
	skip   bool
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: packet-dump <packet-file> <description-file>")
		os.Exit(1)
	}

	pkt, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	rules, err := loadRules(os.Args[2])
	if err != nil {
		panic(err)
	}

	segments := []segment{}

	offset := 0
	for _, rule := range rules {
		bytes, length, err := rule.Matcher.Match(pkt[offset:])
		if err != nil {
			panic(err)
		}

		segments = append(segments, segment{
			offset: offset,
			length: length,
			name:   rule.Description,
			bytes:  bytes,
			skip:   rule.Skip,
		})

		offset += length
	}

	fmt.Print(terminalPresenter{}.Present(segments))
}
