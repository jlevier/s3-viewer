package ui

import (
	"fmt"
	"s3-viewer/api"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	bm bucketsModel
)

type bucketsModel struct {
	buckets []*s3.Bucket
}

type getBucketsMsg struct {
	buckets []*s3.Bucket
	err     error
}

func (m *Model) BucketsInit() tea.Cmd {
	bm = bucketsModel{}

	b, err := api.GetBuckets(m.session)

	if err != nil {
		return func() tea.Msg {
			return getBucketsMsg{nil, err}
		}
	}

	return func() tea.Msg {
		return getBucketsMsg{b, nil}
	}
}

func (m *Model) GetBucketsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case getBucketsMsg:
		if msg.err != nil {
			panic(msg.err) //TODO do something actually meaningful here
		}

		bm.buckets = msg.buckets
	}

	return m, nil
}

func (m *Model) GetBucketsView() string {
	var b strings.Builder

	if bm.buckets != nil {
		for _, bucket := range bm.buckets {
			fmt.Fprintf(&b, "%s\n", *bucket.Name)
		}
		return b.String()
	}

	return "YOU ARE NOW IN THE Buckets VIEW"
}
