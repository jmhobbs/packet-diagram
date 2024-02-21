package presenter

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jmhobbs/packet-diagram/internal/grok"
	"github.com/muesli/termenv"
)

type Terminal struct{}

func (t Terminal) Present(allSegments []grok.Segment, color bool) string {
	if color {
		lipgloss.SetColorProfile(termenv.TrueColor)
	} else {
		lipgloss.SetColorProfile(termenv.Ascii)
	}

	// flatten our skips into blocks
	segments := []grok.Segment{}
	for _, s := range allSegments {
		if len(segments) > 0 && s.Skip && segments[len(segments)-1].Skip {
			continue

		}
		segments = append(segments, s)
	}

	var builder strings.Builder

	offsetStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	nullStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	asciiPrintableStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	asciiWhitespaceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	nonAsciiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

	// find the longest name, as that is the only variable portion of the output
	longestName := 0
	for _, segment := range segments {
		if len(segment.Name) > longestName {
			longestName = len(segment.Name)
		}
	}

	// pre-build the heavy and light box lines for the name segments
	nameBoxLine := ""
	nameBoxLineHeavy := ""
	for i := 0; i < longestName; i++ {
		nameBoxLine += "┄"
		nameBoxLineHeavy += "━"
	}

	// build a format string to create name strings of the correct length and alignment
	nameFmtString := fmt.Sprintf("%%-%ds", longestName)

	// build and align the text for a skipped segment
	skipText := fmt.Sprintf("┃%s┃", lipgloss.NewStyle().
		Width(81+longestName).
		Align(lipgloss.Center).Render("--- skipped ---"))

	// build the lines for between segments, normal, skip above and skip below
	separator := "┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄" + nameBoxLine + "┄┨"
	skipJoinBottomSseparator := "┠┄┄┄┄┄┄┄┄┬┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┬┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┬┄┄┄┄┄┄┄┄┬┄┄┄┄┄┄┄┄┬┄" + nameBoxLine + "┄┨"
	skipJoinTopSseparator := "┠┄┄┄┄┄┄┄┄┴┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┴┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┴┄┄┄┄┄┄┄┄┴┄┄┄┄┄┄┄┄┴┄" + nameBoxLine + "┄┨"

	// write the header, differentiating if the first segment is skipped
	if segments[0].Skip {
		builder.WriteString("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	} else {
		builder.WriteString("┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━")
	}
	builder.WriteString(nameBoxLineHeavy)
	builder.WriteString("━┓\n")

	for j, segment := range segments {
		if segment.Skip {
			builder.WriteString(skipText)
			builder.WriteRune('\n')
			if j+1 < len(segments) {
				builder.WriteString(skipJoinBottomSseparator)
				builder.WriteRune('\n')
			}
			continue
		}

		lines := int(math.Ceil(float64(len(segment.Bytes)) / 16.0))

		for line := 0; line < lines; line++ {
			lineBytes := segment.Bytes[line*16 : min(segment.Length, (line+1)*16)]

			// build up strings for hex encoded bytes, and the printable representations
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

			// only print the name for the first line of a segment
			name := fmt.Sprintf(nameFmtString, "")
			if line == 0 {
				name = fmt.Sprintf(nameFmtString, segment.Name)
			}

			fmt.Fprintf(
				&builder,
				"┃%s│ %s ┊ %s │%s┊%s│ %s ┃\n",
				offsetStyle.Render(fmt.Sprintf("%08d", segment.Offset+(16*line))),
				strings.Join(bytes[:8], " "),
				strings.Join(bytes[8:], " "),
				strings.Join(printables[:8], ""),
				strings.Join(printables[8:], ""),
				name,
			)
		}

		if j < len(segments)-1 {
			if segments[j+1].Skip {
				builder.WriteString(skipJoinTopSseparator)
			} else {
				builder.WriteString(separator)
			}
			builder.WriteRune('\n')
		}
	}

	if segments[len(segments)-1].Skip {
		builder.WriteString("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	} else {
		builder.WriteString("┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━")
	}
	builder.WriteString(nameBoxLineHeavy)
	builder.WriteString("━┛\n")

	return builder.String()
}
