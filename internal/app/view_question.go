package app

import "charm.land/lipgloss/v2"

type Question struct {
	question string
}

func (q Question) View(prompt string, width int) string {
	styleQuestion := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAF0E6"))
	styleQuestionPrompt := lipgloss.NewStyle().Padding(1, 0)
	styleQuestionPrompt = styleQuestionPrompt.Width(width - styleQuestionPrompt.GetHorizontalFrameSize() - 1)

	return styleQuestionPrompt.Render(prompt + " " + styleQuestion.Render(q.question))
}
