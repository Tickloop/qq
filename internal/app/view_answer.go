package app

import (
	glamour "charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"
)

var r *glamour.TermRenderer

func buildRenderer(width int) {
	s := styles.DarkStyleConfig
	s.Document.BlockPrefix = ""
	r, _ = glamour.NewTermRenderer(
		glamour.WithStyles(s),
		glamour.WithWordWrap(width),
	)
}

func styleAnswer(answer string) string {
	if r == nil {
		return answer
	}
	out, err := r.Render(answer)
	if err != nil {
		return answer
	}
	return out
}

func styleError(err string) string {
	style := lipgloss.NewStyle()
	style = style.Foreground(lipgloss.BrightRed)
	style = style.Bold(true)
	return style.Render(err)
}
