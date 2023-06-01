package control

import (
	"s3-viewer/ui/buckets"
	"s3-viewer/ui/creds"
	"s3-viewer/ui/files"
	"s3-viewer/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	uiModel *types.UiModel
)

// Use of empty model here just to satisfy the tea contract for Init, View, and Model.
// The reason we are not passing types.UiModel and implementing the methods on that struct
// is because we would like to keep the ui pages in separate packages and that would require
// that they reference the tea model, thus creating a circular reference error since those packages
// must also be referenced here by the control to pass off page functionality.
type Model struct{}

func (m Model) Init() tea.Cmd {
	if uiModel == nil {
		uiModel = types.GetInitialModel()
	}

	switch uiModel.GetCurrentPage() {
	case types.Buckets:
		return buckets.Init(uiModel)
	case types.Files:
		return files.Init(uiModel)
	default:
		return creds.Init(uiModel)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" {
			return m, tea.Quit
		}

	case types.ChangeCurrentPageMsg:
		switch uiModel.GetCurrentPage() {
		case types.Buckets:
			return m, buckets.Init(uiModel)
		case types.Files:
			return m, files.Init(uiModel)
		default:
			return m, creds.Init(uiModel)
		}
	}

	switch uiModel.GetCurrentPage() {
	case types.Buckets:
		return m, buckets.Update(uiModel, msg)
	case types.Files:
		return m, files.Update(uiModel, msg)
	default:
		return m, creds.Update(uiModel, msg)
	}
}

func (m Model) View() string {
	switch uiModel.GetCurrentPage() {
	case types.Buckets:
		return buckets.View(uiModel)
	case types.Files:
		return files.View(uiModel)
	default:
		return creds.View(uiModel)
	}
}
