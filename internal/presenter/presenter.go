package presenter

import "github.com/jmhobbs/packet-diagram/internal/grok"

type Presenter interface {
	Present([]grok.Segment) string
}
