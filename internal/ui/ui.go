package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	tickRateMs = 16 // The rate of UI updates.
	padding    = 2  // The size of UI padding; applied in all dimensions.
	maxWidth   = 80 // The max character width of the UI.
)

type tickMsg time.Time

type UI struct {
	spinner  spinner.Model
	progress []*Progress
}

func NewUI(p ...*Progress) *UI {
	return &UI{
		spinner:  spinner.New(spinner.WithSpinner(spinner.MiniDot)),
		progress: p,
	}
}

func (m UI) Init() tea.Cmd {
	for _, p := range m.progress {
		p.Init()
	}

	return tea.Batch(tick(), m.spinner.Tick)
}

func (m UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		var completed int

		// Iterate over the progress bars
		for _, p := range m.progress {
			if *p.percentage >= 1 {
				completed++
				continue
			}

			p.Update(msg)
		}

		if completed == len(m.progress) {
			return m, tea.Quit
		}

		cmds = append(cmds, tick())

	case spinner.TickMsg:
		s, cmd := m.spinner.Update(msg)
		m.spinner = s

		cmds = append(cmds, cmd)

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m UI) View() string {
	var out, bars string

	var completed int

	// Iterate over the progress bars
	for _, p := range m.progress {
		if *p.percentage >= 1 {
			completed++
			continue
		}

		// Render the row
		bars += fmt.Sprintf("%s %s\n", strings.Repeat(" ", padding), p.View())
	}

	if completed == len(m.progress) {
		return ""
	}

	out += fmt.Sprintf("%s [%d/%d] Completing action...\n", m.spinner.View(), completed, len(m.progress))

	return out + bars
}

func tick() tea.Cmd {
	return tea.Tick(tickRateMs*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
