package ui

import (
	"s3-viewer/api"

	"github.com/aws/aws-sdk-go/aws/session"
)

type CurrentPage string

const (
	Creds   CurrentPage = "creds"
	Buckets             = "buckets"
	Main                = "main"
)

type Model struct {
	currentPage    CurrentPage
	session        *session.Session
	loadingMessage string
	errorMessage   string
}

func InitialModel() Model {
	ch := make(chan *api.SessionResponse)
	go api.GetSession(ch)
	resp := <-ch

	if resp.Err != nil {
		return Model{
			currentPage: Creds,
			session:     nil,
		}
	}

	return Model{
		currentPage: Buckets,
		session:     resp.Session,
	}
}
