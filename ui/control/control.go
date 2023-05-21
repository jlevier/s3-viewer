package control

import (
	"s3-viewer/api"
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

func getInitialModel() *types.UiModel {
	ch := make(chan *api.SessionResponse)
	go api.GetSession(ch)
	resp := <-ch

	if resp.Err != nil {
		return &types.UiModel{
			CurrentPage: types.Creds,
			Session:     nil,
		}
	}

	return &types.UiModel{
		CurrentPage: types.Buckets,
		Session:     resp.Session,
	}
}

func (m Model) Init() tea.Cmd {
	if uiModel == nil {
		uiModel = getInitialModel()
	}

	// TODO - need to refactor how the current page is being changed because this Init method
	// is only being called once.  Should probably create a changeCurrentPageMsg that can bubble up to the Update
	// method here
	// cmds := make([]tea.Cmd, 3)
	// cmds = append(cmds, buckets.Init(uiModel))
	// cmds = append(cmds, files.Init(uiModel))
	// cmds = append(cmds, creds.Init(uiModel))

	// return tea.Batch(cmds...)

	switch uiModel.CurrentPage {
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
		uiModel.CurrentPage = msg.CurrentPage
		uiModel.CurrentBucket = msg.CurrentBucket

		switch uiModel.CurrentPage {
		case types.Buckets:
			return m, buckets.Init(uiModel)
		case types.Files:
			return m, files.Init(uiModel)
		default:
			return m, creds.Init(uiModel)
		}
	}

	switch uiModel.CurrentPage {
	case types.Buckets:
		return m, buckets.Update(uiModel, msg)
	case types.Files:
		return m, files.Update(uiModel, msg)
	default:
		return m, creds.Update(uiModel, msg)
	}
}

func (m Model) View() string {
	switch uiModel.CurrentPage {
	case types.Buckets:
		return buckets.View(uiModel)
	case types.Files:
		return files.View(uiModel)
	default:
		return creds.View(uiModel)
	}
}
