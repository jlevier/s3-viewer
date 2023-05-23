package buckets

import (
	"os"
	"s3-viewer/api"
	"s3-viewer/ui"
	"s3-viewer/ui/components/table"
	"s3-viewer/ui/types"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	model     *bucketsModel
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))
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
			r = append(r, table.Row{*b.Name, b.CreationDate.String()})
		}
		model.table.SetData(r)

		// case tea.KeyMsg:
		// 	switch msg.String() {
		// 	case "esc":
		// 		if model.table.Focused() {
		// 			model.table.Blur()
		// 		} else {
		// 			model.table.Focus()
		// 		}
		// 	case "enter":
		// 		cmds = append(cmds, m.SetCurrentPage(types.Files, &model.table.SelectedRow()[0]))
		// 	}

		// 	var cmd tea.Cmd
		// 	model.table, cmd = model.table.Update(msg)
		// 	cmds = append(cmds, cmd)
	}

	// Default commands
	var sc tea.Cmd
	model.spinner, sc = model.spinner.Update(msg)
	cmds = append(cmds, sc)

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

		//model.table.SetHeight(height - 5)

		p := lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			baseStyle.Render(model.table.View()),
		)

		return docStyle.Render(p)
	}

	return "YOU ARE NOW IN THE Buckets VIEW"
}
