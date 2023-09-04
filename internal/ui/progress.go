package ui

import (
	"math"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Progress struct {
	progress   progress.Model
	percentage *float64

	label string
}

func NewProgressBar(label string, p *float64, opts ...progress.Option) *Progress {
	return &Progress{progress: progress.New(opts...), label: label, percentage: p}
}

func (m Progress) Init() tea.Cmd {
	return nil
}

// NOTE: The 'Progress' implementation is purely a view; percentage updates need
// to be managed externally via updating its internal 'percentage' pointer.
func (m Progress) Update(_ tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	return m, nil
}

func (m Progress) View() string {
	return m.label + m.progress.ViewAs(math.Min(math.Max(*m.percentage, 0), 1))
}
