package buckets

import (
	"fmt"
	"s3-viewer/api"
	"s3-viewer/ui"
	"s3-viewer/ui/types"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	model bucketsModel
)

type bucketsModel struct {
	buckets   []*s3.Bucket
	spinner   spinner.Model
	isLoading bool
}

type getBucketsMsg struct {
	buckets []*s3.Bucket
	err     error
}

func Init(m *types.UiModel) tea.Cmd {
	model = bucketsModel{spinner: ui.GetSpinner(), isLoading: true}

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
	switch msg := msg.(type) {
	case getBucketsMsg:
		if msg.err != nil {
			panic(msg.err) //TODO do something actually meaningful here
		}

		model.buckets = msg.buckets
	}

	// Default commands
	defaultCmds := make([]tea.Cmd, 0)
	var sc tea.Cmd
	model.spinner, sc = model.spinner.Update(msg)
	defaultCmds = append(defaultCmds, sc)

	return tea.Batch(defaultCmds...)
}

func View(m *types.UiModel) string {
	if model.isLoading {
		return ui.GetLoadingDialog("Loading Buckets", model.spinner)
	}

	var b strings.Builder

	if model.buckets != nil {
		for _, bucket := range model.buckets {
			fmt.Fprintf(&b, "%s\n", *bucket.Name)
		}
		return b.String()
	}

	return "YOU ARE NOW IN THE Buckets VIEW"
}
