package render

import (
	"github.com/charmbracelet/glamour"
)

func RenderMarkdown(markdown string) string {
	out, err := glamour.Render(markdown, "dark")
	if err != nil {
		// default: no markdown render
		return markdown
	}
	return out
}
