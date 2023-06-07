package table

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

type Column struct {
	Width int
	Name  string
}

type Row []string

type Model struct {
	columns             []Column
	data                []Row
	highlightedRowIndex int
	firstVisibleRow     int
	footerInfo          string

	// Loading
	spinner   spinner.Model
	isLoading bool

	// Filter
	hasFiltering    bool
	currentFilter   string
	isFilterVisible bool
	filterInput     textinput.Model

	// Paging
	hasNextPage      bool
	currentPageIndex int
}

type FilterAppliedMsg struct {
	Filter string
}

type NextPageMsg struct {
	CurrentPageIndex int
}

type PrevPageMsg struct {
	CurrentPageIndex int
}
