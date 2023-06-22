package files

import (
	"fmt"
	"os"
	"s3-viewer/api"
	"s3-viewer/ui/components/dialog"
	"s3-viewer/ui/components/help"
	"s3-viewer/ui/components/icons"
	spin "s3-viewer/ui/components/spinner"
	"s3-viewer/ui/components/table"
	"s3-viewer/ui/types"
	"s3-viewer/ui/utils"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	model *filesModel
)

type filesModel struct {
	directories        []string
	files              []*s3.Object
	spinner            spinner.Model
	isLoading          bool
	table              *table.Model
	continuationTokens []*string // Used for current, next, previous page
}

type getFilesMsg struct {
	objects *s3.ListObjectsV2Output
	err     error
}

func initTable() *table.Model {
	columns := []table.Column{
		{Name: "", Width: 3}, // Icon column
		{Name: "Key", Width: 50},
		{Name: "Size", Width: 15},
		{Name: "Last Modified", Width: 35},
		{Name: "Owner", Width: 35},
	}

	return table.New(columns, true)
}

func getFileRow(f *s3.Object) table.Row {
	var owner string = ""
	if f.Owner != nil {
		owner = *f.Owner.DisplayName
	}

	return table.Row{
		icons.GetIcon(*f.Key),
		*f.Key,
		utils.GetFriendlyByteDisplay(*f.Size),
		f.LastModified.Format(time.DateTime),
		owner,
	}
}

func createGetFilesMsg(m *types.UiModel, path, filter string, continuationToken *string) func() tea.Msg {
	return func() tea.Msg {
		o, err := api.GetObjects(m.Session, m.GetCurrentBucket(), path, filter, continuationToken)
		model.isLoading = false
		if err != nil {
			return getFilesMsg{nil, err}
		}

		m.SetCurrentPath(path)
		return getFilesMsg{o, nil}
	}
}

func Init(m *types.UiModel) tea.Cmd {
	model = &filesModel{
		spinner:            spin.GetSpinner(),
		isLoading:          true,
		table:              initTable(),
		continuationTokens: make([]*string, 0),
	}

	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, model.spinner.Tick)
	cmds = append(cmds, model.table.Init())

	cmds = append(cmds, func() tea.Msg {
		o, err := api.GetObjects(m.Session, m.GetCurrentBucket(), "/", "", nil)
		model.isLoading = false
		if err != nil {
			return getFilesMsg{nil, err}
		}

		return getFilesMsg{o, nil}
	})

	return tea.Batch(cmds...)
}

func Update(m *types.UiModel, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case getFilesMsg:
		handleGetFilesMsg(m, msg)

	case table.FilterAppliedMsg:
		handleFilterAppliedMsg(m, msg, &cmds)

	case table.NextPageMsg:
		handleNextPageMsg(m, msg, &cmds)

	case table.PrevPageMsg:
		handlePrevPageMsg(m, msg, &cmds)

	case tea.KeyMsg:
		// Filter is visible so allow the table to handle this command and hide the filter
		if model.table.IsFilterVisible() {
			var cmd tea.Cmd
			model.table, cmd = model.table.Update(msg)
			cmds = append(cmds, cmd)
			break
		}

		switch msg.String() {
		case "esc":
			// If a filter is currently applied, do nothing and let table clear the filter
			if model.table.GetCurrentFilter() != "" {
				break
			}

			handleEscKeyMsg(m, msg, &cmds)

		case "enter":
			handleEnterKeyMsg(m, msg, &cmds)
		}
	}

	if model.isLoading {
		var sc tea.Cmd
		model.spinner, sc = model.spinner.Update(msg)
		cmds = append(cmds, sc)
	}

	if model.table.IsFilterVisible() {
		// KeyMsg is handled above so you only want to forward any other type like BlinkMsg
		if _, ok := msg.(tea.KeyMsg); !ok {
			var fc tea.Cmd
			model.table, fc = model.table.Update(msg)
			cmds = append(cmds, fc)
		}
	} else {
		// otherwise pass all messages down to the table
		var fc tea.Cmd
		model.table, fc = model.table.Update(msg)
		cmds = append(cmds, fc)
	}

	return tea.Batch(cmds...)
}

func View(m *types.UiModel) string {
	if model.isLoading {
		return dialog.GetLoadingDialog(fmt.Sprintf("Loading Bucket %s", m.GetCurrentBucket()), model.spinner)
	}

	if model.directories != nil || model.files != nil {
		// Get terminal size and place dialog in the center
		docStyle := lipgloss.NewStyle()
		width, height, _ := term.GetSize(int(os.Stdout.Fd()))

		if width > 0 {
			docStyle = docStyle.MaxWidth(width)
		}
		if height > 0 {
			docStyle = docStyle.MaxHeight(height)
		}

		final := lipgloss.JoinVertical(
			lipgloss.Center,
			model.table.View(),
			help.GetFilesHelp(model.table.IsFilterVisible(), model.table.GetCurrentFilter()))

		p := lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			final,
		)

		return docStyle.Render(p)
	}

	return "YOU ARE NOW IN THE Buckets VIEW"
}
