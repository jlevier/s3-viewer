package files

import (
	"fmt"
	"os"
	"s3-viewer/api"
	"s3-viewer/ui"
	"s3-viewer/ui/components/icons"
	"s3-viewer/ui/components/table"
	"s3-viewer/ui/types"

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
	directories []string
	files       []*s3.Object
	spinner     spinner.Model
	isLoading   bool
	table       *table.Model
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

	return table.New(columns)
}

func getFileRow(f *s3.Object) table.Row {
	var owner string = ""
	if f.Owner != nil {
		owner = *f.Owner.DisplayName
	}

	return table.Row{icons.GetIcon(*f.Key), *f.Key, fmt.Sprint(*f.Size), f.LastModified.String(), owner}
}

func Init(m *types.UiModel) tea.Cmd {
	// if model != nil {
	// 	return nil
	// }

	model = &filesModel{
		spinner:   ui.GetSpinner(),
		isLoading: true,
		table:     initTable(),
	}

	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, model.spinner.Tick)

	cmds = append(cmds, func() tea.Msg {
		o, err := api.GetObjects(m.Session, m.GetCurrentBucket(), "/")
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
		if msg.err != nil {
			panic(msg.err) //TODO do something actually meaningful here
		}

		model.directories = make([]string, len(msg.objects.CommonPrefixes))
		for i, p := range msg.objects.CommonPrefixes {
			model.directories[i] = *p.Prefix
		}
		model.files = msg.objects.Contents

		r := make([]table.Row, 0)
		for _, d := range model.directories {
			r = append(r, table.Row{icons.GetDirectoryIcon(), d, "", "", ""})
		}
		if model.files != nil {
			for _, f := range model.files {
				r = append(r, getFileRow(f))
			}
		}
		model.table.SetData(r)
		model.table.SetFooterInfo(fmt.Sprintf("%s/%s", m.GetCurrentBucket(), m.GetCurrentPath()))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmds = append(cmds, m.SetCurrentPage(types.Buckets, nil))
		}

		var cmd tea.Cmd
		model.table, cmd = model.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Default commands
	var sc tea.Cmd
	model.spinner, sc = model.spinner.Update(msg)
	cmds = append(cmds, sc)

	return tea.Batch(cmds...)
}

func View(m *types.UiModel) string {
	if model.isLoading {
		return ui.GetLoadingDialog(fmt.Sprintf("Loading Bucket %s", m.GetCurrentBucket()), model.spinner)
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

		p := lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			model.table.View(),
		)

		return docStyle.Render(p)
	}

	return "YOU ARE NOW IN THE Buckets VIEW"
}
