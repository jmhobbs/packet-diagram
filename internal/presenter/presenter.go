package presenter

import "github.com/jmhobbs/packet-diagram/internal/grok"

type Config struct {
	ForceColor bool
}

type Presenter interface {
	Present([]grok.Segment, Config) string
}
