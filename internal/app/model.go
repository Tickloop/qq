package app

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	glamour "charm.land/glamour/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/tickloop/qq/internal/agent"
	"github.com/tickloop/qq/internal/config"
	"github.com/tickloop/qq/internal/inference"

	"golang.org/x/term"
)


func styleQuestion(question string) string {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
	}

	styleQuestion := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAF0E6"))
	stylePrompt := lipgloss.NewStyle().Foreground(lipgloss.Color("#00bb00"))

	styleQuestionPrompt := lipgloss.NewStyle().Padding(1, 0)
	styleQuestionPrompt = styleQuestionPrompt.Width(w - styleQuestionPrompt.GetHorizontalFrameSize() - 1)
	
	return styleQuestionPrompt.Render(stylePrompt.Render("\033[1m›\033[0m ") + styleQuestion.Render(question))
}

func styleAnswer(answer string) string {
	out, err := glamour.Render(answer, "dark")
	if err != nil {
		return answer
	}
	return out
}


var providerConverseFnMap = map[string]func(c context.Context, q, m string) (string, error){
	"openrouter": func(c context.Context, q, m string) (string, error) { return inference.OpenRouterConverse(c, q, m) },
	"bedrock":    func(c context.Context, q, m string) (string, error) { return inference.AWSConverse(c, q, m) },
}

type QAModel struct {
	Question string
	Answer string
    Error error
	IsLoading bool
    args config.CLIArgs 
    spinner Spinner
}

type answerMsg struct {
    answer string
    err error
}

func NewQAModel(args config.CLIArgs) QAModel {
    return QAModel{
        Question: args.Question,
        Answer: "",
        Error: nil,
        IsLoading: true,
        args: args,
        spinner: NewSpinner(), 
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
    hldr, ok := providerConverseFnMap[m.args.Provider]
    if !ok {
        return answerMsg{ answer: "", err: fmt.Errorf("ERR: Provider not found") }
    }

    ctx := context.Background()
    questionWithSearchResults := addSearchResultToQuestion(m.args.Question)
    answer, err := hldr(ctx, questionWithSearchResults, m.args.ModelId)
    if err != nil {
        return answerMsg{ answer: "", err: fmt.Errorf("ERR: %v", err) }
    }
    return answerMsg{ answer: answer, err: err }
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
        if m.IsLoading {
            return m, cmd
        }
        return m, nil
    case answerMsg:
        m.IsLoading = false
        m.Answer = msg.answer
        m.Error = msg.err
        return m, tea.Quit
    }
    return m, nil
}

func (m QAModel) View() tea.View {
    questionView := styleQuestion(m.Question)
    answerView := ""
    switch {
    case m.IsLoading:
        answerView = m.spinner.View()
    case m.Error != nil:
        answerView = styleAnswer(m.Error.Error())
    default:
        answerView = styleAnswer(m.Answer)
    }
    return tea.NewView(
        lipgloss.JoinVertical(
            lipgloss.Left,
            questionView,
            answerView,
        ),
    )
}

