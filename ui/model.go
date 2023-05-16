package ui

import (
	"s3-viewer/s3"

	"github.com/aws/aws-sdk-go/aws/session"
)

type CurrentPage string

const (
	Creds CurrentPage = "creds"
	Main              = "main"
)

type Model struct {
	currentPage CurrentPage
	session     *session.Session
}

func InitialModel() Model {
	ch := make(chan *s3.SessionResponse)
	go s3.GetSession(ch)
	resp := <-ch

	if resp.Err != nil {
		return Model{
			currentPage: Creds,
			session:     nil,
		}
	}

	return Model{
		currentPage: Main,
		session:     resp.Session,
	}
}
