package curl

import (
	"io"

	"github.com/fatih/color"
)

// Multiple formats being exported for later reuse
var (
	HeaderFormat      = color.New(color.Bold, color.FgBlue)
	TitleFormat       = color.New(color.Bold, color.Underline, color.FgRed)
	StandardFormat    = color.New(color.FgWhite)
	CommandFormat     = color.New(color.BgBlack, color.FgHiGreen)
	ExplanationFormat = color.New(color.Italic, color.FgWhite)
)

// Header writes a text using the HeaderFormat to the w io.Writer
func Header(w io.Writer, text string) {
	HeaderFormat.Fprintln(w, text)
}

// Title writes a text using the TitleFormat to the w io.Writer
func Title(w io.Writer, text string) {
	TitleFormat.Fprintln(w, text)
}

// Standard writes a text using the StandardFormat to the w io.Writer
func Standard(w io.Writer, text string) {
	StandardFormat.Fprintln(w, text)
}

// Command writes a text using the CommandFormat to the w io.Writer
func Command(w io.Writer, text string) {
	CommandFormat.Fprintln(w, text)
}

// Explanation writes a text using the ExplanationFormat to the w io.Writer
func Explanation(w io.Writer, text string) {
	ExplanationFormat.Fprintln(w, text)
}

func init() {
	color.NoColor = false
}
