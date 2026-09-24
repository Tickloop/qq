package app

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
type TickMsg struct {}
type Spinner struct { position int }

func NewSpinner() Spinner {
    return Spinner{ position: 0 }
}

func styleFrame(f string) string {
    style := lipgloss.NewStyle()
    style = style.Foreground(lipgloss.Color("#00ff00"))
    return style.Render(f)
}

func (s Spinner) Tick() tea.Msg {
    time.Sleep(100 * time.Millisecond)
    return TickMsg{}
}

func (s Spinner) Init() tea.Cmd {
    return nil
}

func (s Spinner) Update(msg tea.Msg) (Spinner, tea.Cmd) {
    switch msg.(type) {
    case TickMsg:
        s.position = ( s.position + 1 ) % len(frames)
        return s, s.Tick
    }
    return s, nil
}

func (s Spinner) View() string {
    return styleFrame(frames[s.position])
}
