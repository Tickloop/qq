package app

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"context"
	"fmt"
	"github.com/tickloop/qq/internal/agent"
	"github.com/tickloop/qq/internal/config"
	"github.com/tickloop/qq/internal/inference"
)

type QAModel struct {
	question   Question
	answer     string
	err        error
	isLoading  bool
	args       config.CLIArgs
	spinner    Spinner
	windowSize windowSize
}

type answerMsg struct {
	answer string
	err    error
}

type windowSize struct {
	w int
	h int
}

func NewQAModel(args config.CLIArgs) QAModel {
	return QAModel{
		question:   Question{question: args.Question},
		answer:     "",
		err:        nil,
		isLoading:  true,
		args:       args,
		spinner:    NewSpinner(),
		windowSize: windowSize{w: 80, h: 24},
	}
}

func addSearchResultToQuestion(question string) string {
	questionWithSearchResults := fmt.Sprintf("Question: %v\n", question)
	searchResults, err := agent.WebSearchFirecrawl(question)
	if err != nil {
		searchResults = "ERR: Search results failed to load"
	} else {
		questionWithSearchResults = fmt.Sprintf("Search Results: %v\n\nQuestion: %v\n", searchResults, question)
	}
	return questionWithSearchResults
}

func (m QAModel) fetchAnswer() tea.Msg {
	hldr := inference.ProviderConverseFnMap[m.args.Provider]
	questionWithSearchResults := addSearchResultToQuestion(m.args.Question)

	ctx := context.Background()
	answer, err := hldr(ctx, questionWithSearchResults, m.args.ModelId)
	if err != nil {
		return answerMsg{answer: "", err: fmt.Errorf("ERR: %v", err)}
	}
	return answerMsg{answer: answer, err: err}
}

func (m QAModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchAnswer)
}

func (m QAModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case TickMsg:
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner
		if m.isLoading {
			return m, cmd
		}
		return m, nil
	case answerMsg:
		m.isLoading = false
		m.answer = msg.answer
		m.err = msg.err
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.windowSize.w = msg.Width
		m.windowSize.h = msg.Height
		buildRenderer(msg.Width)
		return m, nil
	}
	return m, nil
}

func (m QAModel) View() tea.View {
	// question section
	promptView := promptView()
	if m.isLoading {
		promptView = m.spinner.View(m.isLoading)
	}
	questionView := m.question.View(promptView, m.windowSize.w)

	// answer section
	answerView := ""
	switch {
	case m.isLoading:
		answerView = ""
	case m.err != nil:
		answerView = styleError(m.err.Error())
	case m.answer != "":
		answerView = styleAnswer(m.answer)
	}

	// compose
	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			questionView,
			answerView,
		),
	)
}
