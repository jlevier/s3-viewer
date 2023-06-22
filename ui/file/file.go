package file

import (
	"fmt"
	"s3-viewer/ui/components/dialog"
	spin "s3-viewer/ui/components/spinner"
	"s3-viewer/ui/types"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	model *fileModel
)

type fileModel struct {
	spinner   spinner.Model
	isLoading bool
}

func Init(m *types.UiModel) tea.Cmd {
	// if model != nil {
	// 	return nil
	// }

	model = &fileModel{
		spinner:   spin.GetSpinner(),
		isLoading: true,
	}

	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, model.spinner.Tick)

	return tea.Batch(cmds...)
}

func Update(m *types.UiModel, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	if model.isLoading {
		var sc tea.Cmd
		model.spinner, sc = model.spinner.Update(msg)
		cmds = append(cmds, sc)
	}

	return tea.Batch(cmds...)
}

func View(m *types.UiModel) string {
	if model.isLoading {
		return dialog.GetLoadingDialog(fmt.Sprintf("Loading File %s", m.GetCurrentFile()), model.spinner)
	}

	return "You are now in the file"
}
