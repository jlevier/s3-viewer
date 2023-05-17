package ui

import (
	"fmt"
	"os"
	"s3-viewer/s3"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	viewModel = initialModel()
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 50

		switch i {
		case 0:
			t.Placeholder = "Key"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Secret"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

type validateCredsMsg struct {
	sess *session.Session
	err  error
}

func validateCreds(key, secret string) tea.Cmd {
	return func() tea.Msg {
		sess, err := s3.GetSessionFromInput(key, secret)
		return validateCredsMsg{sess, err}
	}
}

func (m *Model) CredsInit() tea.Cmd {
	return textinput.Blink
}

func (m *Model) GetCredsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()

		switch k {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == "enter" && viewModel.focusIndex == len(viewModel.inputs) {
				m.loadingMessage = "Validating..."
				m.errorMessage = ""
				return m, validateCreds(viewModel.inputs[0].Value(), viewModel.inputs[1].Value())
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				viewModel.focusIndex--
			} else {
				viewModel.focusIndex++
			}

			if viewModel.focusIndex > len(viewModel.inputs) {
				viewModel.focusIndex = 0
			} else if viewModel.focusIndex < 0 {
				viewModel.focusIndex = len(viewModel.inputs)
			}

			cmds := make([]tea.Cmd, len(viewModel.inputs))
			for i := 0; i <= len(viewModel.inputs)-1; i++ {
				if i == viewModel.focusIndex {
					// Set focused state
					//cmds[i] = viewModel.inputs[i].Focus()
					cmds = append(cmds, viewModel.inputs[i].Focus())
					viewModel.inputs[i].PromptStyle = focusedStyle
					viewModel.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				viewModel.inputs[i].Blur()
				viewModel.inputs[i].PromptStyle = noStyle
				viewModel.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}

	case validateCredsMsg:
		m.loadingMessage = ""
		if msg.err != nil {
			m.errorMessage = msg.err.Error()
		} else {
			m.currentPage = Main
			m.session = msg.sess
		}
		return m, nil
	}

	cmd := viewModel.updateInputs(msg)

	return m, cmd
}

func (m *Model) GetCredsView() string {
	var b strings.Builder

	h1 := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Seems you don't have any cached credentials.")
	h2 := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Enter your AWS key and secret.")
	header := lipgloss.JoinVertical(lipgloss.Center, h1, h2)
	fmt.Fprintf(&b, "%s\n\n", header)

	for i := range viewModel.inputs {
		// Provide padding on the front of the text boxes
		b.WriteString(" ")
		b.WriteString(viewModel.inputs[i].View())
		if i < len(viewModel.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := buttonStyle.Render("Submit")
	if viewModel.focusIndex == len(viewModel.inputs) {
		button = activeButtonStyle.Render("Submit")
	}
	buttonAligned := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(button)
	fmt.Fprintf(&b, "\n\n%s", buttonAligned)

	if m.loadingMessage != "" {
		fmt.Fprintf(&b, "\n\n%s", m.loadingMessage)
	} else if m.errorMessage != "" {
		fmt.Fprintf(&b, "\n\n%s", m.errorMessage)
	} else {
		b.WriteString("\n\n")
	}

	// Get terminal size and place dialog in the center
	docStyle := lipgloss.NewStyle()
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	if width > 0 {
		docStyle = docStyle.MaxWidth(width)
	}
	if height > 0 {
		docStyle = docStyle.MaxHeight(height)
	}

	p := lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(b.String()),
		lipgloss.WithWhitespaceChars("Ш#"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}))

	return docStyle.Render(p)
}
