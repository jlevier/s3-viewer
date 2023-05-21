package files

import (
	"fmt"
	"s3-viewer/api"
	"s3-viewer/ui/types"

	"github.com/aws/aws-sdk-go/service/s3"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	model filesModel
)

type filesModel struct {
	files []*s3.Object
	//spinner   spinner.Model
	isLoading bool
	//table     table.Model
}

type getFilesMsg struct {
	objects *s3.ListObjectsV2Output
	err     error
}

func Init(m *types.UiModel) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	//cmds = append(cmds, model.spinner.Tick)

	cmds = append(cmds, func() tea.Msg {
		o, err := api.GetObjects(m.Session, m.CurrentBucket)
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

		model.files = msg.objects.Contents

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmds = append(cmds, func() tea.Msg {
				return types.ChangeCurrentPageMsg{
					CurrentPage: types.Buckets,
				}
			})
		}
	}

	return tea.Batch(cmds...)
}

func View(m *types.UiModel) string {
	return fmt.Sprintf("In the files view for %s\n%s", m.CurrentBucket, model.files)
}
