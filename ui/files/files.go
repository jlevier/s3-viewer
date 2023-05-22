package files

import (
	"s3-viewer/api"
	"s3-viewer/ui/list"
	"s3-viewer/ui/types"

	"github.com/aws/aws-sdk-go/service/s3"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	model filesModel
)

type filesModel struct {
	directories []string
	files       []*s3.Object
	//spinner   spinner.Model
	isLoading bool
	//table     table.Model
	listModel *list.Model
}

type getFilesMsg struct {
	objects *s3.ListObjectsV2Output
	err     error
}

func Init(m *types.UiModel) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	//cmds = append(cmds, model.spinner.Tick)

	model.listModel = list.NewModel()

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

	cmds = append(cmds, list.Update(msg, model.listModel))

	switch msg := msg.(type) {
	case getFilesMsg:
		if msg.err != nil {
			panic(msg.err) //TODO do something actually meaningful here
		}

		model.directories = make([]string, len(msg.objects.CommonPrefixes))
		for _, p := range msg.objects.CommonPrefixes {
			model.directories = append(model.directories, *p.Prefix)
		}
		model.files = msg.objects.Contents

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmds = append(cmds, m.SetCurrentPage(types.Buckets, nil))
		}
	}

	return tea.Batch(cmds...)
}

func View(m *types.UiModel) string {
	// s := fmt.Sprintf("%s", model.directories)
	// return fmt.Sprintf("In the files view for %s\n%s\n%s", s, m.GetCurrentBucket(), model.files)
	return list.View(model.listModel)
}
