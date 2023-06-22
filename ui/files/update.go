package files

import (
	"fmt"
	"s3-viewer/ui/components/icons"
	"s3-viewer/ui/components/table"
	"s3-viewer/ui/types"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func handleGetFilesMsg(m *types.UiModel, msg getFilesMsg) {
	if msg.err != nil {
		panic(msg.err) //TODO do something actually meaningful here
	}

	model.directories = make([]string, len(msg.objects.CommonPrefixes))
	for i, p := range msg.objects.CommonPrefixes {
		model.directories[i] = *p.Prefix
	}
	model.files = msg.objects.Contents

	if msg.objects.NextContinuationToken != nil {
		model.continuationTokens = append(model.continuationTokens, msg.objects.NextContinuationToken)
		model.table.SetHasNextPage(true)
	} else {
		model.table.SetHasNextPage(false)
	}

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
}

func handleFilterAppliedMsg(m *types.UiModel, msg table.FilterAppliedMsg, cmds *[]tea.Cmd) {
	*cmds = append(*cmds, createGetFilesMsg(m, m.GetCurrentPath(), msg.Filter, nil))
}

func handleNextPageMsg(m *types.UiModel, msg table.NextPageMsg, cmds *[]tea.Cmd) {
	*cmds = append(
		*cmds,
		createGetFilesMsg(
			m,
			m.GetCurrentPath(),
			model.table.GetCurrentFilter(),
			model.continuationTokens[msg.CurrentPageIndex]))
}

func handlePrevPageMsg(m *types.UiModel, msg table.PrevPageMsg, cmds *[]tea.Cmd) {
	// pop off the last continuation token
	model.continuationTokens = model.continuationTokens[:len(model.continuationTokens)-1]
	var ct *string
	if msg.CurrentPageIndex > 0 {
		ct = model.continuationTokens[msg.CurrentPageIndex-1]
	}

	*cmds = append(
		*cmds,
		createGetFilesMsg(
			m,
			m.GetCurrentPath(),
			model.table.GetCurrentFilter(),
			ct))
}

func handleEscKeyMsg(m *types.UiModel, msg tea.KeyMsg, cmds *[]tea.Cmd) {
	// At the root of the current bucket - return to buckets list
	if m.GetCurrentPath() == "" {
		*cmds = append(*cmds, m.SetCurrentPage(types.Buckets, nil))
	} else {
		// drilled into a folder inside of a bucket - go one folder up
		p := strings.Split(m.GetCurrentPath(), "/")
		cp := fmt.Sprintf("%s%s", strings.Join(p[:len(p)-2], "/"), "/")
		if cp == "/" {
			cp = ""
		}

		*cmds = append(*cmds, createGetFilesMsg(m, cp, "", nil))
	}
}

func handleEnterKeyMsg(m *types.UiModel, msg tea.KeyMsg, cmds *[]tea.Cmd) {
	r := model.table.GetHighlightedRow()
	*cmds = append(*cmds, createGetFilesMsg(m, (*r)[1], "", nil))
}
