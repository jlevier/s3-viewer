package ui

import (
	"fmt"
	"s3-viewer/api"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	bm bucketsModel
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

func (m *Model) BucketsInit() tea.Cmd {
	bm = bucketsModel{spinner: getSpinner(), isLoading: true}

	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, bm.spinner.Tick)

	cmds = append(cmds, func() tea.Msg {
		b, err := api.GetBuckets(m.session)
		bm.isLoading = false
		if err != nil {
			return getBucketsMsg{nil, err}
		}
		return getBucketsMsg{b, nil}
	})

	return tea.Batch(cmds...)
}

func (m *Model) GetBucketsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case getBucketsMsg:
		if msg.err != nil {
			panic(msg.err) //TODO do something actually meaningful here
		}

		bm.buckets = msg.buckets
	}

	// Default commands
	defaultCmds := make([]tea.Cmd, 0)
	var sc tea.Cmd
	bm.spinner, sc = bm.spinner.Update(msg)
	defaultCmds = append(defaultCmds, sc)

	return m, tea.Batch(defaultCmds...)
}

func (m *Model) GetBucketsView() string {
	if bm.isLoading {
		return getLoadingDialog("Loading Buckets", bm.spinner)
	}

	var b strings.Builder

	if bm.buckets != nil {
		for _, bucket := range bm.buckets {
			fmt.Fprintf(&b, "%s\n", *bucket.Name)
		}
		return b.String()
	}

	return "YOU ARE NOW IN THE Buckets VIEW"
}
