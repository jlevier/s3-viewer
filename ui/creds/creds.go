package creds

import (
	"fmt"
	"os"
	"s3-viewer/api"
	"s3-viewer/ui"
	"s3-viewer/ui/types"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/charmbracelet/bubbles/spinner"
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
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4754")).Width(55).Padding(2)
	dialogHeaderStyle   = lipgloss.NewStyle().Width(50).Align(lipgloss.Center)
	buttonAlignedStyle  = lipgloss.NewStyle().Width(50).Align(lipgloss.Center)

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

	model = initialModel()
)

type credsModel struct {
	focusIndex     int
	inputs         []textinput.Model
	cursorMode     textinput.CursorMode
	spinner        spinner.Model
	loadingMessage string
	errorMessage   string
}

func initialModel() credsModel {
	m := credsModel{
		inputs:  make([]textinput.Model, 2),
		spinner: ui.GetSpinner(),
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

func (m *credsModel) updateInputs(msg tea.Msg) tea.Cmd {
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
		sess, err := api.GetSessionFromInput(key, secret)
		return validateCredsMsg{sess, err}
	}
}

func Init(m *types.UiModel) tea.Cmd {
	return tea.Batch(textinput.Blink, model.spinner.Tick)
}

func Update(m *types.UiModel, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()

		switch k {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == "enter" && model.focusIndex == len(model.inputs) {
				model.loadingMessage = "Validating..."
				model.errorMessage = ""
				return validateCreds(model.inputs[0].Value(), model.inputs[1].Value())
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				model.focusIndex--
			} else {
				model.focusIndex++
			}

			if model.focusIndex > len(model.inputs) {
				model.focusIndex = 0
			} else if model.focusIndex < 0 {
				model.focusIndex = len(model.inputs)
			}

			cmds := make([]tea.Cmd, len(model.inputs))
			for i := 0; i <= len(model.inputs)-1; i++ {
				if i == model.focusIndex {
					// Set focused state
					//cmds[i] = model.inputs[i].Focus()
					cmds = append(cmds, model.inputs[i].Focus())
					model.inputs[i].PromptStyle = focusedStyle
					model.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				model.inputs[i].Blur()
				model.inputs[i].PromptStyle = noStyle
				model.inputs[i].TextStyle = noStyle
			}

			return tea.Batch(cmds...)
		}

	case validateCredsMsg:
		model.loadingMessage = ""
		if msg.err != nil {
			model.errorMessage = fmt.Sprintf("\u274C %s", msg.err.Error())
		} else {
			m.Session = msg.sess
			//TODO need to make this return the cmd (declare defaultCmds below higher up and append here)
			m.SetCurrentPage(types.Buckets, nil)
		}
		return nil
	}

	// Default commands
	defaultCmds := make([]tea.Cmd, 0)
	defaultCmds = append(defaultCmds, model.updateInputs(msg))
	var sc tea.Cmd
	model.spinner, sc = model.spinner.Update(msg)
	defaultCmds = append(defaultCmds, sc)

	return tea.Batch(defaultCmds...)
}

func View(m *types.UiModel) string {
	var b strings.Builder

	h1 := dialogHeaderStyle.Render("Seems you don't have any cached credentials.")
	h2 := dialogHeaderStyle.Render("Enter your AWS key and secret:")
	header := lipgloss.JoinVertical(lipgloss.Center, h1, h2)
	fmt.Fprintf(&b, "%s\n\n", header)

	for i := range model.inputs {
		// Provide padding on the front of the text boxes
		b.WriteString(" ")
		b.WriteString(model.inputs[i].View())
		if i < len(model.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := buttonStyle.Render("Submit")
	if model.focusIndex == len(model.inputs) {
		button = activeButtonStyle.Render("Submit")
	}
	fmt.Fprintf(&b, "\n\n%s", buttonAlignedStyle.Render(button))

	if model.loadingMessage != "" {
		fmt.Fprintf(&b, "\n\n%s%s", model.spinner.View(), model.loadingMessage)
	} else if model.errorMessage != "" {
		fmt.Fprintf(&b, "\n\n%s", errorStyle.Render(model.errorMessage))
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
		ui.DialogBoxStyle.Render(b.String()),
		lipgloss.WithWhitespaceChars("Ш#"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}))

	return docStyle.Render(p)
}
