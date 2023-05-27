package buckets

import (
	"os"
	"s3-viewer/api"
	"s3-viewer/ui"
	"s3-viewer/ui/components/table"
	"s3-viewer/ui/types"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	model *bucketsModel

	iconStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E87C3C"))
)

type bucketsModel struct {
	buckets   []*s3.Bucket
	spinner   spinner.Model
	isLoading bool
	table     *table.Model
}

type getBucketsMsg struct {
	buckets []*s3.Bucket
	err     error
}

func initTable() *table.Model {
	columns := []table.Column{
		{Name: "", Width: 3}, // Icon column
		{Name: "Bucket", Width: 50},
		{Name: "Creation Date", Width: 35},
	}

	return table.New(columns)
}

func Init(m *types.UiModel) tea.Cmd {
	if model != nil {
		return nil
	}

	model = &bucketsModel{
		spinner:   ui.GetSpinner(),
		isLoading: true,
		table:     initTable(),
	}

	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, model.spinner.Tick)

	cmds = append(cmds, func() tea.Msg {
		b, err := api.GetBuckets(m.Session)
		model.isLoading = false
		if err != nil {
			return getBucketsMsg{nil, err}
		}
		return getBucketsMsg{b, nil}
	})

	return tea.Batch(cmds...)
}

func Update(m *types.UiModel, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case getBucketsMsg:
		if msg.err != nil {
			panic(msg.err) //TODO do something actually meaningful here
		}

		model.buckets = msg.buckets
		r := make([]table.Row, 0)
		for _, b := range model.buckets {
			r = append(r, table.Row{iconStyle.Render("\ue703"), *b.Name, b.CreationDate.Format(time.DateTime)})
		}
		model.table.SetData(r)

	case tea.KeyMsg:
		switch msg.String() {
		// 	case "esc":
		// 		if model.table.Focused() {
		// 			model.table.Blur()
		// 		} else {
		// 			model.table.Focus()
		// 		}
		case "enter":
			r := model.table.GetHighlightedRow()
			if r != nil {
				cmds = append(cmds, m.SetCurrentPage(types.Files, &(*r)[1]))
			}
		}

		var cmd tea.Cmd
		model.table, cmd = model.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	if model.isLoading {
		var sc tea.Cmd
		model.spinner, sc = model.spinner.Update(msg)
		cmds = append(cmds, sc)
	}

	return tea.Batch(cmds...)
}

func View(m *types.UiModel) string {
	if model.isLoading {
		return ui.GetLoadingDialog("Loading Buckets", model.spinner)
	}

	if model.buckets != nil {
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
			model.table.View(),
		)

		return docStyle.Render(p)
	}

	return "YOU ARE NOW IN THE Buckets VIEW"
}
