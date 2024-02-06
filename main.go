package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type segment struct {
	offset int
	length int
	name   string
	bytes  []byte
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
		})

		offset += length
	}

	longestName := 0
	for _, segment := range segments {
		if len(segment.name) > longestName {
			longestName = len(segment.name)
		}
	}

	nameFmtString := fmt.Sprintf("%%-%ds", longestName)

	fmt.Print("┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━")
	for i := 0; i < longestName; i++ {
		fmt.Print("━")
	}
	fmt.Println("━┓")

	separator := "┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄"
	for i := 0; i < longestName; i++ {
		separator += "┄"
	}
	separator += "┄┨"

	for j, segment := range segments {
		lines := int(math.Ceil(float64(len(segment.bytes)) / 16.0))

		for line := 0; line < lines; line++ {
			lineBytes := segment.bytes[line*16 : min(segment.length, (line+1)*16)]

			bytes := []string{}
			printables := []string{}
			for _, b := range lineBytes {
				bytes = append(bytes, fmt.Sprintf("%02x", b))
				// TODO: Color
				if b >= 32 && b <= 126 {
					printables = append(printables, string(b))
				} else {
					// TODO: Alternative printables for special characters
					printables = append(printables, "⋄")
				}
			}

			// add blanks for incomplete lines
			for i := len(lineBytes); i < 16; i++ {
				bytes = append(bytes, "  ")
				printables = append(printables, " ")
			}

			name := fmt.Sprintf(nameFmtString, "")

			if line == 0 {
				name = fmt.Sprintf(nameFmtString, segment.name)
			}

			fmt.Printf(
				"┃%08d│ %s ┊ %s │%s┊%s│ %s ┃\n",
				segment.offset+(16*line),
				strings.Join(bytes[:8], " "),
				strings.Join(bytes[8:], " "),
				strings.Join(printables[:8], ""),
				strings.Join(printables[8:], ""),
				name,
			)
		}

		if j < len(segments)-1 {
			fmt.Println(separator)
		}
	}

	fmt.Print("┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━")
	for i := 0; i < longestName; i++ {
		fmt.Print("━")
	}
	fmt.Println("━┛")

}
