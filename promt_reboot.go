package main

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type rebootModel struct {
	cursor  int
	choice  chan int
	choices []string
}

func (m rebootModel) Init() tea.Cmd {
	return nil
}

func (m rebootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "esc":
			os.Exit(1)
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice <- m.cursor
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}

		}
	}

	return m, nil
}

func (m rebootModel) View() string {
	s := strings.Builder{}

	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n\n(esc to quit)\n")

	return s.String()
}
