package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coffeebeats/gdenv/internal/ui"
)

func main() {
	var v1, v2, v3 float64
	p1, p2, p3 := &v1, &v2, &v3

	var bars []*ui.Progress

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)

		switch i {
		case 0:
			bars = append(bars, ui.NewProgressBar("Downloading "+string(i), &v1))

			go func() {
				for *p1 < 1 {
					*p1 += float64(rand.Intn(4)) / float64(100)
					time.Sleep(5 * time.Millisecond) //nolint:wsl,gomnd
				}

				wg.Done()
			}()
		case 1:
			bars = append(bars, ui.NewProgressBar("Downloading "+string(i), &v2))

			go func() {
				for *p2 < 1 {
					*p2 += float64(rand.Intn(5)) / float64(100)
					time.Sleep(25 * time.Millisecond)
				}

				wg.Done()
			}()
		case 2:
			bars = append(bars, ui.NewProgressBar("Downloading "+string(i), &v3))

			go func() {
				for *p3 < 1 {
					*p3 += float64(rand.Intn(5)) / float64(100)
					time.Sleep(50 * time.Millisecond)
				}

				wg.Done()
			}()
		}

	}

	fmt.Println("Start")

	p := tea.NewProgram(ui.NewUI(bars...))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	fmt.Println("DONE")

	tea.ClearScreen()
}

// package main

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/charmbracelet/bubbles/spinner"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// type tickMsg time.Time

// var (
// 	// Available spinners
// 	spinners = []spinner.Spinner{
// 		spinner.Line,
// 		spinner.Dot,
// 		spinner.MiniDot,
// 		spinner.Jump,
// 		spinner.Pulse,
// 		spinner.Points,
// 		spinner.Globe,
// 		spinner.Moon,
// 		spinner.Monkey,
// 	}

// 	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
// 	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
// 	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
// )

// func main() {
// 	m := model{}
// 	m.resetSpinner()

// 	if _, err := tea.NewProgram(m).Run(); err != nil {
// 		fmt.Println("could not run program:", err)
// 		os.Exit(1)
// 	}
// }

// type model struct {
// 	index   int
// 	spinner spinner.Model
// }

// func (m model) Init() tea.Cmd {
// 	return tea.Batch(tick(), m.spinner.Tick)
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

// 	var cmds []tea.Cmd

// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q", "esc":
// 			return m, tea.Quit
// 		case "h", "left":
// 			m.index--
// 			if m.index < 0 {
// 				m.index = len(spinners) - 1
// 			}
// 			m.resetSpinner()
// 			cmds = append(cmds, m.spinner.Tick)
// 		case "l", "right":
// 			m.index++
// 			if m.index >= len(spinners) {
// 				m.index = 0
// 			}
// 			m.resetSpinner()
// 			cmds = append(cmds, m.spinner.Tick)
// 		default:
// 			return m, nil
// 		}
// 	case tickMsg:
// 		cmds = append(cmds, tick())
// 	case spinner.TickMsg:
// 		s, cmd := m.spinner.Update(msg)
// 		m.spinner = s
// 		cmds = append(cmds, cmd)
// 	}

// 	return m, tea.Batch(cmds...)
// }

// func (m *model) resetSpinner() {
// 	m.spinner = spinner.New()
// 	m.spinner.Style = spinnerStyle
// 	m.spinner.Spinner = spinners[m.index]
// }

// func (m model) View() (s string) {
// 	var gap string
// 	switch m.index {
// 	case 1:
// 		gap = ""
// 	default:
// 		gap = " "
// 	}

// 	s += fmt.Sprintf("\n %s%s%s\n\n", m.spinner.View(), gap, textStyle("Spinning..."))
// 	s += helpStyle("h/l, ←/→: change spinner • q: exit\n")
// 	return
// }

// func tick() tea.Cmd {
// 	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
// 		return tickMsg(t)
// 	})
// }
