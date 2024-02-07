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

func (t terminalPresenter) Present(allSegments []segment) string {

	// flatten our skips into blocks
	segments := []segment{}
	for _, s := range allSegments {
		if len(segments) > 0 && s.skip && segments[len(segments)-1].skip {
			continue

		}
		segments = append(segments, s)
	}

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

	descriptionLine := ""
	descriptionLineHeavy := ""
	for i := 0; i < longestName; i++ {
		descriptionLine += "┄"
		descriptionLineHeavy += "━"
	}

	nameFmtString := fmt.Sprintf("%%-%ds", longestName)

	if segments[0].skip {
		builder.WriteString("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	} else {
		builder.WriteString("┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━")
	}
	builder.WriteString(descriptionLineHeavy)
	builder.WriteString("━┓\n")

	skipText := fmt.Sprintf("┃%s┃", lipgloss.NewStyle().
		Width(81+longestName).
		Align(lipgloss.Center).Render("--- skipped ---"))

	separator := "┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄" + descriptionLine + "┄┨"
	skipJoinBottomSseparator := "┠┄┄┄┄┄┄┄┄┬┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┬┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┬┄┄┄┄┄┄┄┄┬┄┄┄┄┄┄┄┄┬┄" + descriptionLine + "┄┨"
	skipJoinTopSseparator := "┠┄┄┄┄┄┄┄┄┴┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┴┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┴┄┄┄┄┄┄┄┄┴┄┄┄┄┄┄┄┄┴┄" + descriptionLine + "┄┨"

	for j, segment := range segments {
		if segment.skip {
			builder.WriteString(skipText)
			builder.WriteRune('\n')
			if j+1 < len(segments) {
				builder.WriteString(skipJoinBottomSseparator)
				builder.WriteRune('\n')
			}
			continue
		}

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
			if segments[j+1].skip {
				builder.WriteString(skipJoinTopSseparator)
			} else {
				builder.WriteString(separator)
			}
			builder.WriteRune('\n')
		}
	}

	if segments[len(segments)-1].skip {
		builder.WriteString("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	} else {
		builder.WriteString("┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━")
	}
	builder.WriteString(descriptionLineHeavy)
	builder.WriteString("━┛\n")

	return builder.String()
}
