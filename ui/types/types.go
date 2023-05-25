package types

import (
	"s3-viewer/api"

	"github.com/aws/aws-sdk-go/aws/session"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	Creds   CurrentPage = "creds"
	Buckets             = "buckets"
	Files               = "files"
)

type CurrentPage string

// Used to change the current page.  Use this msg instead of just setting directly the current page on the
// UiModel because it is necessary for the main Update() method in control to call the Init() function of the new page.
type ChangeCurrentPageMsg struct {
	CurrentPage   CurrentPage
	CurrentBucket string
}

// This is the main model used for the overall UI and for
// pages to pass information back and forth to each other.
type UiModel struct {
	Session       *session.Session
	currentPage   CurrentPage
	currentBucket string
	currentPath   string
}

func GetInitialModel() *UiModel {
	ch := make(chan *api.SessionResponse)
	go api.GetSession(ch)
	resp := <-ch

	if resp.Err != nil {
		m := &UiModel{
			currentPage: Creds,
			Session:     nil,
		}

		return m
	}

	m := &UiModel{
		currentPage: Buckets,
		Session:     resp.Session,
	}

	return m
}

func (m *UiModel) GetCurrentPage() CurrentPage {
	return m.currentPage
}

func (m *UiModel) GetCurrentBucket() string {
	return m.currentBucket
}

func (m *UiModel) GetCurrentPath() string {
	return m.currentPath
}

func (m *UiModel) SetCurrentPage(currentPage CurrentPage, currentBucket *string) tea.Cmd {
	if currentBucket != nil {
		m.currentBucket = *currentBucket
	} else {
		m.currentBucket = ""
	}

	return func() tea.Msg {
		m.currentPage = currentPage
		return ChangeCurrentPageMsg{
			CurrentPage:   m.currentPage,
			CurrentBucket: m.currentBucket,
		}
	}
}
