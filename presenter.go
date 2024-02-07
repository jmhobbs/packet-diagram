package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Presenter interface {
	Present([]byte) string
}

type terminalPresenter struct {
}

func (t terminalPresenter) Present(segments []segment) string {
	var builder strings.Builder

	nullStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))
	asciiPrintableStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00aaff"))
	asciiWhitespaceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff33"))
	nonAsciiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF22"))

	longestName := 0
	for _, segment := range segments {
		if len(segment.name) > longestName {
			longestName = len(segment.name)
		}
	}

	nameFmtString := fmt.Sprintf("%%-%ds", longestName)

	builder.WriteString("┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━")
	for i := 0; i < longestName; i++ {
		builder.WriteRune('━')
	}
	builder.WriteString("━┓\n")

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
				if b == 0 {
					bytes = append(bytes, nullStyle.Render(fmt.Sprintf("%02x", b)))
					printables = append(printables, nullStyle.Render("⋄"))
				} else if b >= 32 && b <= 126 {
					if b == 0x20 || b == 0x09 || b == 0x0D || b == 0x0A {
						bytes = append(bytes, asciiWhitespaceStyle.Render(fmt.Sprintf("%02x", b)))
						printables = append(printables, asciiWhitespaceStyle.Render(" "))
					} else {
						bytes = append(bytes, asciiPrintableStyle.Render(fmt.Sprintf("%02x", b)))
						printables = append(printables, asciiPrintableStyle.Render(string(b)))
					}
				} else {
					bytes = append(bytes, nonAsciiStyle.Render(fmt.Sprintf("%02x", b)))
					printables = append(printables, nonAsciiStyle.Render("x"))
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

			fmt.Fprintf(
				&builder,
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
			builder.WriteString(separator)
			builder.WriteRune('\n')
		}
	}

	builder.WriteString("┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━")
	for i := 0; i < longestName; i++ {
		builder.WriteRune('━')
	}
	builder.WriteString("━┛\n")

	return builder.String()
}
