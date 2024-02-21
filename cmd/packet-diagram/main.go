package main

import (
	"fmt"
	"os"

	"github.com/jmhobbs/packet-diagram/internal/grok"
	"github.com/jmhobbs/packet-diagram/internal/presenter"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: packet-dump <packet-file> <description-file>")
		os.Exit(1)
	}

	pkt, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	rules, err := grok.LoadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	segments := []grok.Segment{}

	offset := 0
	for _, rule := range rules {
		segment, err := rule.Match(pkt, offset)
		if err != nil {
			panic(err)
		}
		segments = append(segments, segment)
		offset += segment.Length
	}

	fmt.Print(presenter.Terminal{}.Present(segments, true))
}
