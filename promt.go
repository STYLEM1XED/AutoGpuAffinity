package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type model struct {
	textInput textinput.Model
	err       error
	title     string
	callback  chan int
}

func initialModel(title, def string, callback chan int) model {
	ti := textinput.New()
	ti.SetValue(def)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		title:     title,
		textInput: ti,
		err:       nil,
		callback:  callback,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			os.Exit(1)
			return m, tea.Quit
		case tea.KeyEnter:
			intVar, err := strconv.Atoi(m.textInput.Value())
			if err == nil {
				m.callback <- intVar
				return m, tea.Quit
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n",
		m.title,
		m.textInput.View(),
		"(esc to quit)",
	)
}
