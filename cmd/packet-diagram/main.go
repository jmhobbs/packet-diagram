package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jmhobbs/packet-diagram/internal/grok"
	"github.com/jmhobbs/packet-diagram/internal/presenter"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

func main() {
	fs := ff.NewFlagSet("packet-diagram [options] <packet-file> <grok-file>")

	var (
		forceColor = fs.Bool('c', "force-color", "Force color output, even if not connected to a terminal")
	)

	err := ff.Parse(fs, os.Args[1:])
	switch {
	case errors.Is(err, ff.ErrHelp):
		fmt.Fprintf(os.Stderr, "%s\n", ffhelp.Flags(fs))
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	args := fs.GetArgs()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "error: expected 2 arguments, got %d\n\n", len(args))
		fmt.Fprintf(os.Stderr, "%s\n", ffhelp.Flags(fs))
		os.Exit(1)
	}

	pkt, err := os.ReadFile(args[0])
	if err != nil {
		panic(err)
	}

	rules, err := grok.LoadFile(args[1])
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

	fmt.Print(presenter.Terminal{}.Present(segments, presenter.Config{ForceColor: *forceColor}))
}
